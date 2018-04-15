package v1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

const (
	CRDPlural   string = "databases"
	CRDGroup    string = "postgresql.org"
	CRDVersion  string = "v1"
	FullCRDName string = CRDPlural + "." + CRDGroup
)

// Create the CRD resource, ignore error if it already exists
//func CreateCRD(clientset apiextcs.Interface) error {
//	crd := &apiextv1beta1.CustomResourceDefinition{
//		ObjectMeta: meta_v1.ObjectMeta{Name: FullCRDName},
//		Spec: apiextv1beta1.CustomResourceDefinitionSpec{
//			Group:   CRDGroup,
//			Version: CRDVersion,
//			Scope:   apiextv1beta1.NamespaceScoped,
//			Names: apiextv1beta1.CustomResourceDefinitionNames{
//				Plural: CRDPlural,
//				Kind:   reflect.TypeOf(Postgres{}).Name(),
//			},
//		},
//	}
//
//	_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
//	if err != nil && apierrors.IsAlreadyExists(err) {
//		return nil
//	}
//	return err
//	// Note the original apiextensions example adds logic to wait for creation and exception handling
//}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Definition of our CRD Postgres class
type Postgres struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               PostgresConfig `json:"spec"`
	Status             PostgresStatus `json:"status,omitempty"`
}

type PostgresConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}
type PostgresSpec struct {
	Foo string `json:"foo"`
	Bar bool   `json:"bar"`
	Baz int    `json:"baz,omitempty"`
}

type PostgresStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []Postgres `json:"items"`
}

func NewClient(cfg *rest.Config) (*rest.RESTClient, *runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}
	config := *cfg
	config.GroupVersion = &SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{
		CodecFactory: serializer.NewCodecFactory(scheme)}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, nil, err
	}
	return client, scheme, nil
}
