package api

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	amconfig "github.com/prometheus/alertmanager/config"
	"github.com/weaveworks/common/user"

	"github.com/cortexproject/cortex/pkg/configs"
	"github.com/cortexproject/cortex/pkg/configs/db"
	"github.com/cortexproject/cortex/pkg/util"
)

type Config struct {
	Notifications NotificationsConfig `yaml:"notifications"`
}

type NotificationsConfig struct {
	DisableEmail     bool `yaml:"allow_email"`
	DisablePagerDuty bool `yaml:"allow_pagerduty"`
	DisablePushover  bool `yaml:"allow_pushover"`
	DisableSlack     bool `yaml:"allow_slack"`
	DisableOpsGenie  bool `yaml:"allow_opsgenie"`
	DisableWebHook   bool `yaml:"allow_webhook"`
	DisableVictorOps bool `yaml:"allow_victorops"`
	DisableWeChat    bool `yaml:"allow_wechat"`
	// Hipchat has reached end of life and is no longer available
}

// RegisterFlagsWithPrefix adds the flags required to config this to the given FlagSet
func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	f.BoolVar(&cfg.Notifications.DisableEmail, "configs-api.notifications.disable-email", false, "Disable Email notifications for Alertmanager.")
	f.BoolVar(&cfg.Notifications.DisablePagerDuty, "configs-api.notifications.disable-pagerduty", false, "Disable PagerDuty notifications for Alertmanager.")
	f.BoolVar(&cfg.Notifications.DisablePushover, "configs-api.notifications.disable-pushover", false, "Disable Pushover notifications for Alertmanager.")
	f.BoolVar(&cfg.Notifications.DisableSlack, "configs-api.notifications.disable-slack", false, "Disable Slack notifications for Alertmanager.")
	f.BoolVar(&cfg.Notifications.DisableOpsGenie, "configs-api.notifications.disable-opsgenie", false, "Disable OpsGenie notifications for Alertmanager.")
	f.BoolVar(&cfg.Notifications.DisableWebHook, "configs-api.notifications.disable-webhook", false, "Disable WebHook notifications for Alertmanager.")
	f.BoolVar(&cfg.Notifications.DisableVictorOps, "configs-api.notifications.disable-victorops", false, "Disable VictorOps notifications for Alertmanager.")
	f.BoolVar(&cfg.Notifications.DisableWeChat, "configs-api.notifications.disable-wechat", false, "Disable WeChat notifications for Alertmanager.")
}

// API implements the configs api.
type API struct {
	http.Handler
	db  db.DB
	cfg Config
}

// New creates a new API
func New(database db.DB, cfg Config) *API {
	a := &API{
		db:  database,
		cfg: cfg,
	}
	r := mux.NewRouter()
	a.RegisterRoutes(r)
	a.Handler = r
	return a
}

func (a *API) admin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, `
<!doctype html>
<html>
	<head><title>configs :: configuration service</title></head>
	<body>
		<h1>configs :: configuration service</h1>
	</body>
</html>
`)
}

// RegisterRoutes registers the configs API HTTP routes with the provided Router.
func (a *API) RegisterRoutes(r *mux.Router) {
	for _, route := range []struct {
		name, method, path string
		handler            http.HandlerFunc
	}{
		{"root", "GET", "/", a.admin},
		// Dedicated APIs for updating rules config. In the future, these *must*
		// be used.
		{"get_rules", "GET", "/api/prom/configs/rules", a.getConfig},
		{"set_rules", "POST", "/api/prom/configs/rules", a.setConfig},
		{"get_templates", "GET", "/api/prom/configs/templates", a.getConfig},
		{"set_templates", "POST", "/api/prom/configs/templates", a.setConfig},
		{"get_alertmanager_config", "GET", "/api/prom/configs/alertmanager", a.getConfig},
		{"set_alertmanager_config", "POST", "/api/prom/configs/alertmanager", a.setConfig},
		{"validate_alertmanager_config", "POST", "/api/prom/configs/alertmanager/validate", a.validateAlertmanagerConfig},
		{"deactivate_config", "DELETE", "/api/prom/configs/deactivate", a.deactivateConfig},
		{"restore_config", "POST", "/api/prom/configs/restore", a.restoreConfig},
		// Internal APIs.
		{"private_get_rules", "GET", "/private/api/prom/configs/rules", a.getConfigs},
		{"private_get_alertmanager_config", "GET", "/private/api/prom/configs/alertmanager", a.getConfigs},
	} {
		r.Handle(route.path, route.handler).Methods(route.method).Name(route.name)
	}
}

// getConfig returns the request configuration.
func (a *API) getConfig(w http.ResponseWriter, r *http.Request) {
	userID, _, err := user.ExtractOrgIDFromHTTPRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	logger := util.WithContext(r.Context(), util.Logger)

	cfg, err := a.db.GetConfig(r.Context(), userID)
	if err == sql.ErrNoRows {
		http.Error(w, "No configuration", http.StatusNotFound)
		return
	} else if err != nil {
		// XXX: Untested
		level.Error(logger).Log("msg", "error getting config", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cfg); err != nil {
		// XXX: Untested
		level.Error(logger).Log("msg", "error encoding config", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *API) setConfig(w http.ResponseWriter, r *http.Request) {
	userID, _, err := user.ExtractOrgIDFromHTTPRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	logger := util.WithContext(r.Context(), util.Logger)

	var cfg configs.Config
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		// XXX: Untested
		level.Error(logger).Log("msg", "error decoding json body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validateAlertmanagerConfig(cfg.AlertmanagerConfig, a.cfg.Notifications); err != nil && cfg.AlertmanagerConfig != "" {
		level.Error(logger).Log("msg", "invalid Alertmanager config", "err", err)
		http.Error(w, fmt.Sprintf("Invalid Alertmanager config: %v", err), http.StatusBadRequest)
		return
	}
	if err := validateRulesFiles(cfg); err != nil {
		level.Error(logger).Log("msg", "invalid rules", "err", err)
		http.Error(w, fmt.Sprintf("Invalid rules: %v", err), http.StatusBadRequest)
		return
	}
	if err := validateTemplateFiles(cfg); err != nil {
		level.Error(logger).Log("msg", "invalid templates", "err", err)
		http.Error(w, fmt.Sprintf("Invalid templates: %v", err), http.StatusBadRequest)
		return
	}
	if err := a.db.SetConfig(r.Context(), userID, cfg); err != nil {
		// XXX: Untested
		level.Error(logger).Log("msg", "error storing config", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) validateAlertmanagerConfig(w http.ResponseWriter, r *http.Request) {
	logger := util.WithContext(r.Context(), util.Logger)
	cfg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		level.Error(logger).Log("msg", "error reading request body", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = validateAlertmanagerConfig(string(cfg), a.cfg.Notifications); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		util.WriteJSONResponse(w, map[string]string{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	util.WriteJSONResponse(w, map[string]string{
		"status": "success",
	})
}

func validateAlertmanagerConfig(cfg string, noCfg NotificationsConfig) error {
	amCfg, err := amconfig.Load(cfg)
	if err != nil {
		return err
	}

	for _, recv := range amCfg.Receivers {
		if noCfg.DisableEmail && len(recv.EmailConfigs) > 0 {
			return fmt.Errorf("email notifications are disabled in Cortex yet")
		}
		if noCfg.DisablePagerDuty && len(recv.PagerdutyConfigs) > 0 {
			return fmt.Errorf("pager-duty notifications are disabled in Cortex yet")
		}
		if noCfg.DisablePushover && len(recv.PushoverConfigs) > 0 {
			return fmt.Errorf("pushover notifications are disabled in Cortex yet")
		}
		if noCfg.DisableSlack && len(recv.SlackConfigs) > 0 {
			return fmt.Errorf("slack notifications are disabled in Cortex yet")
		}
		if noCfg.DisableOpsGenie && len(recv.OpsGenieConfigs) > 0 {
			return fmt.Errorf("ops-genie notifications are disabled in Cortex yet")
		}
		if noCfg.DisableWebHook && len(recv.WebhookConfigs) > 0 {
			return fmt.Errorf("web-hook notifications are disabled in Cortex yet")
		}
		if noCfg.DisableVictorOps && len(recv.VictorOpsConfigs) > 0 {
			return fmt.Errorf("victor-ops notifications are disabled in Cortex yet")
		}
		if noCfg.DisableWeChat && len(recv.WechatConfigs) > 0 {
			return fmt.Errorf("we-chat notifications are disabled in Cortex yet")
		}
	}

	return nil
}

func validateRulesFiles(c configs.Config) error {
	_, err := c.RulesConfig.Parse()
	return err
}

func validateTemplateFiles(c configs.Config) error {
	for fn, content := range c.TemplateFiles {
		if _, err := template.New(fn).Parse(content); err != nil {
			return err
		}
	}

	return nil
}

// ConfigsView renders multiple configurations, mapping userID to configs.View.
// Exposed only for tests.
type ConfigsView struct {
	Configs map[string]configs.View `json:"configs"`
}

func (a *API) getConfigs(w http.ResponseWriter, r *http.Request) {
	var cfgs map[string]configs.View
	var cfgErr error
	logger := util.WithContext(r.Context(), util.Logger)
	rawSince := r.FormValue("since")
	if rawSince == "" {
		cfgs, cfgErr = a.db.GetAllConfigs(r.Context())
	} else {
		since, err := strconv.ParseUint(rawSince, 10, 0)
		if err != nil {
			level.Info(logger).Log("msg", "invalid config ID", "err", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		cfgs, cfgErr = a.db.GetConfigs(r.Context(), configs.ID(since))
	}

	if cfgErr != nil {
		// XXX: Untested
		level.Error(logger).Log("msg", "error getting configs", "err", cfgErr)
		http.Error(w, cfgErr.Error(), http.StatusInternalServerError)
		return
	}

	view := ConfigsView{Configs: cfgs}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(view); err != nil {
		// XXX: Untested
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *API) deactivateConfig(w http.ResponseWriter, r *http.Request) {
	userID, _, err := user.ExtractOrgIDFromHTTPRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	logger := util.WithContext(r.Context(), util.Logger)

	if err := a.db.DeactivateConfig(r.Context(), userID); err != nil {
		if err == sql.ErrNoRows {
			level.Info(logger).Log("msg", "deactivate config - no configuration", "userID", userID)
			http.Error(w, "No configuration", http.StatusNotFound)
			return
		}
		level.Error(logger).Log("msg", "error deactivating config", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	level.Info(logger).Log("msg", "config deactivated", "userID", userID)
	w.WriteHeader(http.StatusOK)
}

func (a *API) restoreConfig(w http.ResponseWriter, r *http.Request) {
	userID, _, err := user.ExtractOrgIDFromHTTPRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	logger := util.WithContext(r.Context(), util.Logger)

	if err := a.db.RestoreConfig(r.Context(), userID); err != nil {
		if err == sql.ErrNoRows {
			level.Info(logger).Log("msg", "restore config - no configuration", "userID", userID)
			http.Error(w, "No configuration", http.StatusNotFound)
			return
		}
		level.Error(logger).Log("msg", "error restoring config", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	level.Info(logger).Log("msg", "config restored", "userID", userID)
	w.WriteHeader(http.StatusOK)
}
