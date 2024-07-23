package aliyun

import (
	"context"
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	cred "github.com/aliyun/credentials-go/credentials"
	"github.com/gin-generator/ginctl/package/get"
	"sync"
)

var (
	driver string
	ak     string
	sk     string
	region string
	token  string
	way    string
	once   sync.Once
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) GetCredentials(_ context.Context) (credentials.Credentials, error) {
	return credentials.Credentials{
		AccessKeyID:     ak,
		AccessKeySecret: sk,
		SecurityToken:   token,
	}, nil
}

func NewOssClient() (client *oss.Client) {
	// loading config
	SetDriver()
	SetAk()
	SetSk()
	SetRegion()
	SetToken()

	config := GetCredential()
	once.Do(func() {
		client = oss.NewClient(config)
	})
	return
}

func SetDriver() {
	driver = get.Get("filesystem.driver")
}

func SetAk() {
	ak = get.Get(fmt.Sprintf("filesystem.%s.access_key", driver))
}

func SetSk() {
	sk = get.Get(fmt.Sprintf("filesystem.%s.secret_key", driver))
}

func SetRegion() {
	region = get.String(fmt.Sprintf("filesystem.%s.region", driver))
}

func SetToken() {
	token = get.Get(fmt.Sprintf("filesystem.%s.token", driver), "")
}

func SetWay() {
	way = get.String(fmt.Sprintf("filesystem.%s.credential_way", driver))
}

// GetCredential 获取凭证
func GetCredential() *oss.Config {
	if way == "" {
		SetWay()
	}
	switch way {
	case "region":
		return GetRegionCred()
	case "long-term":
		return GetLongTermCred()
	case "ecs":
		return GetEcsCred()
	case "process":
		return GetOutProcessCred()
	case "ram":
		return GetRamCred()
	default:
		panic("Unsupported way to get credentials.")
	}
}

// GetRegionCred 区域凭证
func GetRegionCred() *oss.Config {
	provider := NewProvider()
	return oss.LoadDefaultConfig().
		WithCredentialsProvider(provider).
		WithRegion(region)
}

// GetLongTermCred 长期凭证
func GetLongTermCred() *oss.Config {
	provider := credentials.NewStaticCredentialsProvider(ak, sk)
	return oss.LoadDefaultConfig().
		WithCredentialsProvider(provider).
		WithRegion(region)
}

// GetEcsCred 指定实例角色获取凭证
func GetEcsCred() *oss.Config {
	provider := credentials.NewEcsRoleCredentialsProvider(func(cpo *credentials.EcsRoleCredentialsProviderOptions) {
		cpo.RamRole = get.Get(fmt.Sprintf("filesystem.%s.ecs_ram_role", driver))
	})
	return oss.LoadDefaultConfig().
		WithCredentialsProvider(provider).
		WithRegion(region)
}

// GetOutProcessCred 获取外部进程长期凭证
func GetOutProcessCred() *oss.Config {
	process := get.Get(fmt.Sprintf("filesystem.%s.process", driver))
	provider := credentials.NewProcessCredentialsProvider(process)
	return oss.LoadDefaultConfig().
		WithCredentialsProvider(provider).
		WithRegion(region)
}

func GetRamCred() *oss.Config {

	SetDriver()
	if ak == "" {
		SetAk()
	}
	if sk == "" {
		SetSk()
	}

	userId := get.Int64(fmt.Sprintf("filesystem.%s.ram.user_id", driver))
	role := get.Get(fmt.Sprintf("filesystem.%s.ram.role", driver))

	config := new(cred.Config).
		// Which type of credential you want
		SetType("ram_role_arn").
		// AccessKeyId of your account
		SetAccessKeyId(ak).
		// AccessKeySecret of your account
		SetAccessKeySecret(sk).
		// Format: acs:ram::USER_Id:role/ROLE_NAME
		SetRoleArn(fmt.Sprintf("acs:ram::%d:role/%s", userId, role)).
		// Role Session Name
		SetRoleSessionName(get.Get(fmt.Sprintf("filesystem.%s.ram.session_name", driver)))

	policy := get.String(fmt.Sprintf("filesystem.%s.ram.policy", driver))
	if policy != "" {
		// Not required, limit the permissions of STS Token
		config.SetPolicy(policy)
	}
	expiration := get.Int(fmt.Sprintf("filesystem.%s.ram.expiration", driver), 3600)
	if expiration > 0 {
		// Not required, limit the Valid time of STS Token
		config.SetRoleSessionExpiration(expiration)
	}

	arnCredential, err := cred.NewCredential(config)
	provider := credentials.CredentialsProviderFunc(func(ctx context.Context) (credentials.Credentials, error) {
		if err != nil {
			return credentials.Credentials{}, err
		}
		info, errs := arnCredential.GetCredential()
		if errs != nil {
			return credentials.Credentials{}, errs
		}
		return credentials.Credentials{
			AccessKeyID:     *info.AccessKeyId,
			AccessKeySecret: *info.AccessKeySecret,
			SecurityToken:   *info.SecurityToken,
		}, nil
	})

	return oss.LoadDefaultConfig().
		WithCredentialsProvider(provider).
		WithRegion(region)
}
