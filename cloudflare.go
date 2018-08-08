package cfip

import (
	"errors"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
)

// CF represents a CloudFlare API with our helper wrappers
type CF struct {
	zone   string
	zoneID string
	client *cloudflare.API
	log    *logrus.Entry
}

// NewClient accepts a CF Key, E-mail and zone name and returns a client
func NewClient(key, email, zone string, opts ...cloudflare.Option) (*CF, error) {
	client, err := cloudflare.New(key, email, opts...)
	if err != nil {
		return nil, err
	}
	logger := logrus.New().WithField("app", "cfipupdate")
	cf := &CF{client: client, zone: zone, log: logger}
	if err := cf.updateZoneID(); err != nil {
		return nil, err
	}
	return cf, nil
}

// Set points the specificed host to the specified IP
func (cf *CF) Set(host, ip string) error {
	logger := cf.log.WithFields(logrus.Fields{
		"ip":   ip,
		"host": host,
	})
	rec := cloudflare.DNSRecord{
		Name: host,
	}
	logger.Debug("looking up existing DNS record")
	recs, err := cf.client.DNSRecords(cf.zoneID, rec)
	if err != nil {
		logger.WithError(err).Error("error fetching existing DNS record")
		return err
	}
	switch len(recs) {
	case 0:
		logger.Info("creating initial DNS record")
		record := cloudflare.DNSRecord{
			Name:    host,
			Content: ip,
			Type:    "A",
			Proxied: false,
		}
		_, err = cf.client.CreateDNSRecord(cf.zoneID, record)
		return err
	case 1:
		record := recs[0]
		if record.Content == ip {
			logger.Info("old ip matched new ip, doing nothing")
			return nil
		}
		record.Content = ip
		logger.WithField("recordID", record.ID).Info("updating DNS record")
		return cf.client.UpdateDNSRecord(cf.zoneID, record.ID, record)
	default:
		return errors.New(">1 records for that host, I want 1")
	}
}

func (cf *CF) updateZoneID() (err error) {
	cf.zoneID, err = cf.client.ZoneIDByName(cf.zone)
	cf.log = cf.log.WithField("zoneID", cf.zoneID)
	cf.log.Print("got zone id")
	return
}
