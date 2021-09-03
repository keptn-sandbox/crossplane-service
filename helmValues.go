package main

type helmValues struct {
	Helmservice struct {
		Image struct {
			Repository string `yaml:"repository"`
			PullPolicy string `yaml:"pullPolicy"`
			Tag        string `yaml:"tag"`
		} `yaml:"image"`
		Service struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"service"`
	} `yaml:"helmservice"`
	Distributor struct {
		StageFilter   string `yaml:"stageFilter"`
		ServiceFilter string `yaml:"serviceFilter"`
		ProjectFilter string `yaml:"projectFilter"`
		Image         struct {
			Repository string `yaml:"repository"`
			PullPolicy string `yaml:"pullPolicy"`
			Tag        string `yaml:"tag"`
		} `yaml:"image"`
	} `yaml:"distributor"`
	RemoteControlPlane struct {
		Enabled bool `yaml:"enabled"`
		API     struct {
			Protocol       string `yaml:"protocol"`
			Hostname       string `yaml:"hostname"`
			APIValidateTLS bool   `yaml:"apiValidateTls"`
			Token          string `yaml:"token"`
		} `yaml:"api"`
	} `yaml:"remoteControlPlane"`
	ImagePullSecrets []interface{} `yaml:"imagePullSecrets"`
	ServiceAccount   struct {
		Create      bool `yaml:"create"`
		Annotations struct {
		} `yaml:"annotations"`
		Name string `yaml:"name"`
	} `yaml:"serviceAccount"`
	PodAnnotations struct {
	} `yaml:"podAnnotations"`
	PodSecurityContext struct {
	} `yaml:"podSecurityContext"`
	SecurityContext struct {
	} `yaml:"securityContext"`
	Resources struct {
		Limits struct {
			CPU    string `yaml:"cpu"`
			Memory string `yaml:"memory"`
		} `yaml:"limits"`
		Requests struct {
			CPU    string `yaml:"cpu"`
			Memory string `yaml:"memory"`
		} `yaml:"requests"`
	} `yaml:"resources"`
	NodeSelector struct {
	} `yaml:"nodeSelector"`
	Tolerations []interface{} `yaml:"tolerations"`
	Affinity    struct {
	} `yaml:"affinity"`
}
