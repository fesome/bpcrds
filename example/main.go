package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	v1 "github.com/fesome/bpcrds/apis/calico/v1"
	calicov1 "github.com/fesome/bpcrds/client/clientset/versioned/typed/calico/v1"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	//config, err := getFormKubeConfig()
	config, err := SelectClusterConfig("develop")
	if err != nil {
		panic(err.Error())
	}

	clientset, err := calicov1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ipPools, err := clientset.IPPools().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d ipPools in the cluster\n", len(ipPools.Items))

	for _, pool := range ipPools.Items {
		poolBytes, err := yaml.Marshal(pool)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(string(poolBytes))
	}

	err = clientset.IPPools().Delete(context.TODO(), "pool-test", metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}

	poolInfo, err := clientset.IPPools().Create(context.TODO(), &v1.IPPool{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pool-test",
			Labels: map[string]string{
				"organization": "test",
			},
		},
		Spec: v1.IPPoolSpec{
			CIDR:             "192.171.0.0/16",
			VXLANMode:        v1.VXLANModeNever,
			IPIPMode:         v1.IPIPModeNever,
			NATOutgoing:      false,
			Disabled:         false,
			DisableBGPExport: false,
			BlockSize:        26,
			NodeSelector:     "cop!=monitor",
		},
	}, metav1.CreateOptions{})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(poolInfo)

	poolBytes, err := yaml.Marshal(poolInfo)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(poolBytes))
}

func getFormKubeConfig() (*rest.Config, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	return config, nil
}

type K8sConfig struct {
	Name  string // 集群名称
	Host  string // 集群地址 https://xxxx:8443
	Token string // Token 存放serviceaccount对应secret的token
	CA    string // CA证书 存放serviceaccount对应secret的ca.crt
}

var (
	developConfig = K8sConfig{
		Name:  "develop",
		Host:  "https://192.168.1.80:6443",
		Token: "eyJhbGciOiJSUzI1NiIsImtpZCI6IkNZZ3F2NXBfaXVmSERHbERUTUJzbGpGbElvUmxGa19VajhJQzh0OEhhancifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjY2Utb3BzIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImNjZS1vcHMtdXNlci10b2tlbi1oN2twdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJjY2Utb3BzLXVzZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI2ZDU5MDVkZS0xZjMwLTQ1YjYtOTg5YS0yYzczMjUxNzBmMTIiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6Y2NlLW9wczpjY2Utb3BzLXVzZXIifQ.PssJ85OfbANek-1oV3yWhzM2Yjtq3FXh_RgKNoGE61chZQlkRyC1YYf2smb6uX9ifWA5ULoASZNLz3glIQ0_T82YGuBSdwxk2IOvr2IJPQF5GwJGcaF3VokI55w_t4RpynsD5q6lhjQWXZmLtm8CTNYTVTSxylVRQacebTDZR9Vxfol7-HTg39ea3FgjBFw4Ho39ItRKGDT4llVePlZIQcCrQvFzd0YNxwJXgQwl18XORtQvTR537G4FUugEVeXHZooYhQIdUmJ08iPpWleDNCQ1s2SShI3wkb41rjH7Usul57kbnjy2JrNPRxVq33i4JOUjE475YbNVzdUCeEeQUg",
		CA:    base64DecodeConfigOrDie("LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvakNDQWVhZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJeU1ETXhPREUzTXpreU1Wb1hEVE15TURNeE5URTNNemt5TVZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBS2NHCmZKeVNzVnNPeXp0Q2JCYTNqVEF0Q0k2YnNUVGVsMWtNUS9WS0FiUTZXL1Bqdi9jcmhQNHV0UlczZXlkMFpSZkcKKy8yWFI3cnNEazJuUytqOWpLcU9nQ3dsbHFIK052MnJScFAzZFlxRkp4em9MNFZhVTh5YURodmw3OFE0d1h1TAoyV0RMUzNiaVNRTGNSZTFvSm0vZVE4RDZYMnNrSExldVdLM0JEUUs5QnNremFFQ1ZzbTN4dnNCS2I2V3EvTlRvClBycXBHWGg5eDhqRHZIQnI5M0pGVEN5WG9DWit1RXg0UVFOT2xuYmlTSVNSU2EvQXBIQ1NuMnhBNXArcnJLWkIKa3hCR0VNRXExZTdLWm0xblZ0bkNlcW9lTXhmeStUaTUrbHRrSElhdjkzczhqTkY4d2hGYndzTGxla0dmTWFJWApZWWtGeGxWaFhhUzN5SnhXOGpNQ0F3RUFBYU5aTUZjd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZQS2t4VTlDbVpqenZQa1J6My9qOUxleEtxQ21NQlVHQTFVZEVRUU8KTUF5Q0NtdDFZbVZ5Ym1WMFpYTXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBSFVhL2VjRlVSSFZpL3VWU2YxMApvYW92MzR6dE1uWHo1VWNscXNxYVBkR01EZ21xTGhSbVQzaHNjWHRGZzJhc2NrdW9ReDdHZzB6SC82MXBITitZCjluSjk3UTczdEdoYldFMUlIZHp6WDhZNEsySXlWaThSejgyRDVGY05Pb1cxbTJRSkxrMVFndUgvUkFxQ3pGKzgKMk5iYm9SNmRRbzFoRHVHS1hjSEFmc1lReVA4WG0zYU90UnNCR3JCaTM1S3I4UmZ4MGlpbWs5aFpPSUE5dkFOaQoxMjdpSmg4RDhXaGh5WGdDZGNUSjJvbmQ5RU81Rk5lTmhYNlI0ZGdUNWdrTG94bVdjSXRFZENvWWZ2czVDMU1JCkl4SHpQZm93V3I5L3F2RkxNcG1XYjk2L20yY0FKV1duWURQcXNINkZnemthaTNEZ2tXTEpuaW9DUUtEeXJJU3AKbEEwPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="),
	}

	prodConfig = K8sConfig{
		Name:  "prod",
		Host:  "https://172.16.20.10:6443",
		Token: "xxxxxxxx",
		CA:    "xxxxxxxx",
	}
)

func base64DecodeConfigOrDie(val string) string {
	buf, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		panic(fmt.Sprintf("base64 decode config %s failed, err: %s", val, err.Error()))
	}
	return string(buf)
}

func SelectClusterConfig(env string) (*rest.Config, error) {
	var c K8sConfig
	switch env { // 多集群支持
	case "develop":
		c = developConfig
	case "prod":
		c = prodConfig
	default:
		log.Printf("环境: %s 不支持", env)
		return nil, fmt.Errorf("环境: %s 不支持", env)
	}

	return &rest.Config{
		Host:            c.Host,
		BearerToken:     c.Token,
		BearerTokenFile: "",
		TLSClientConfig: rest.TLSClientConfig{
			//Insecure: true,  // 设置为true时 不需要CA
			CAData: []byte(c.CA),
		},
	}, nil
}
