package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/Masterminds/sprig/v3"
	v1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	namespace  string
	kubeconfig string
	list       bool

	defaultKubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
)

func init() {
	flag.StringVar(&namespace, "n", "default", "-n default")
	flag.StringVar(&kubeconfig, "kubeconfig", defaultKubeconfig, "-kubeconfig ~/.kube/config or from env 'KUBECONFIG'")
	flag.BoolVar(&list, "l", false, "-l list group/resources/verbs")
}

func main() {
	flag.Parse()
	envKubeconfig := os.Getenv("KUBECONFIG")
	if envKubeconfig != "" {
		kubeconfig = envKubeconfig
	}

	fmt.Printf("namespace '%s', kubeconfig '%s'\n", namespace, kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	if list {
		_, lists, err := clientset.ServerGroupsAndResources()
		if err != nil {
			log.Println(err)
		}
		var allVerbs = Perm{}
		for _, api := range lists {
			for _, resource := range api.APIResources {
				gv, _ := schema.ParseGroupVersion(api.GroupVersion)
				group := gv.Group
				allVerbs.AddGroup(group)
				allVerbs.AddResources(group, resource.Name)
				allVerbs.AddResourceVerbs(group, resource.Name, resource.Verbs...)
			}
		}
		temp.Execute(os.Stdout, map[string]any{"Groups": allVerbs})
		return
	}

	create, err := clientset.AuthorizationV1().SelfSubjectRulesReviews().Create(context.TODO(), &v1.SelfSubjectRulesReview{
		Spec: v1.SelfSubjectRulesReviewSpec{
			Namespace: namespace,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	var p = Perm{}
	for _, rule := range create.Status.ResourceRules {
		p.AddGroup(rule.APIGroups...)
		for _, group := range rule.APIGroups {
			p.AddResources(group, rule.Resources...)
			for _, resource := range rule.Resources {
				p.AddResourceVerbs(group, resource, rule.Verbs...)
			}
		}
	}
	fmt.Println(p)
}

var temp, _ = template.New("").Funcs(sprig.FuncMap()).Parse(`
{{ range $k, $g := .Groups }}
{{ if eq $k "" }}""{{ else }}{{ $k }}{{ end }}

  {{- range $rk, $rv := $g }}
    {{- printf "- %s" $rk | nindent 2 }}
    {{- range $verb, $_ := $rv }}
      {{- printf "- %s" $verb | nindent 4 }}
    {{- end }}
  {{- end }}
{{- end }}
`)

type (
	Group    = string
	Resource = string
	Verb     = string
)

type Perm map[Group]map[Resource]map[Verb]struct{}

func (p Perm) String() string {
	bf := &bytes.Buffer{}
	temp.Execute(bf, map[string]any{"Groups": p})
	return bf.String()
}

func (p Perm) AddResourceVerbs(g string, r string, verbs ...string) {
	fn := func(verb string) {
		p[g][r][verb] = struct{}{}
	}
	for _, verb := range verbs {
		fn(verb)
	}
}

func (p Perm) AddResources(g string, r ...string) {
	fn := func(resource string) {
		m := p[g]
		_, ok := m[resource]
		if !ok {
			m[resource] = map[string]struct{}{}
		}
	}
	for _, s := range r {
		fn(s)
	}
}
func (p Perm) AddGroup(g ...string) {
	fn := func(group string) {
		_, ok := p[group]
		if !ok {
			p[group] = make(map[string]map[string]struct{})
		}
	}
	for _, s := range g {
		fn(s)
	}
}
