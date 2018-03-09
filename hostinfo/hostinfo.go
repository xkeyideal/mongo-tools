package hostinfo

import (
	"time"

	"github.com/xkeyideal/mongo-tools/common/db"

	"gopkg.in/mgo.v2/bson"
)

type HostInfoOutput struct {
	System *System `bson:"system" json:"system"`
	Os     *Os     `bson:"os" json:"os"`
	Extra  *Extra  `bson:"extra" json:"extra"`
	Ok     int     `bson:"ok" json:"ok"`
}

type System struct {
	CurrentTime time.Time `bson:"currentTime" json:"currentTime"`
	Hostname    string    `bson:"hostname" json:"hostname"`
	CpuAddrSize int       `bson:"cpuAddrSize" json:"cpuAddrSize"`
	MemSizeMB   int       `bson:"memSizeMB" json:"memSizeMB"`
	NumCores    int       `bson:"numCores" json:"numCores"`
	CpuArch     string    `bson:"cpuArch" json:"cpuArch"`
	NumaEnabled bool      `bson:"numaEnabled" json:"numaEnabled"`
}

type Os struct {
	Type    string `bson:"type" json:"type"`
	Name    string `bson:"name" json:"name"`
	Version string `bson:"version" json:"version"`
}

type Extra struct {
	VersionString   string `bson:"versionString" json:"versionString"`
	LibcVersion     string `bson:"libcVersion" json:"libcVersion"`
	KernelVersion   string `bson:"kernelVersion" json:"kernelVersion"`
	CpuFrequencyMHz string `bson:"cpuFrequencyMHz" json:"cpuFrequencyMHz"`
	CpuFeatures     string `bson:"cpuFeatures" json:"cpuFeatures"`
	PageSize        int    `bson:"pageSize" json:"pageSize"`
	NumPages        int    `bson:"numPages" json:"numPages"`
	MaxOpenFiles    int    `bson:"maxOpenFiles" json:"maxOpenFiles"`
}

//direct=true
type HostInfo struct {
	// for connecting to the db
	SessionProvider *db.SessionProvider
}

func NewHostInfo(sp *db.SessionProvider) *HostInfo {
	return &HostInfo{
		SessionProvider: sp,
	}
}

//https://docs.mongodb.com/v3.2/reference/command/hostInfo/
func (ds *HostInfo) Run() (*HostInfoOutput, error) {
	session, err := ds.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.SetSocketTimeout(0)

	dest := &HostInfoOutput{}
	err = session.DB("admin").Run(bson.D{{"hostInfo", 1}}, dest)

	return dest, err
}
