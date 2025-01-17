package release

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/google/uuid"
	"github.com/openshift/oc-mirror/v2/internal/pkg/api/v2alpha1"
	clog "github.com/openshift/oc-mirror/v2/internal/pkg/log"
	"github.com/openshift/oc-mirror/v2/internal/pkg/mirror"
	//nolint
)

const (
	SignatureURL    string = "https://mirror.openshift.com/pub/openshift-v4/signatures/openshift/release/"
	SignatureDir    string = "/signatures/"
	ContentType     string = "Content-Type"
	ApplicationJson string = "application/json"
)

var (
	defaultPK = `-----BEGIN PGP PUBLIC KEY BLOCK-----

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
=syRO
-----END PGP PUBLIC KEY BLOCK-----`
)

type CincinnatiSchema struct {
	Log              clog.PluggableLoggerInterface
	Config           *v2alpha1.ImageSetConfiguration
	Opts             mirror.CopyOptions
	Client           Client
	Signature        SignatureInterface
	Fail             bool
	CincinnatiParams CincinnatiParams
}

type CincinnatiParams struct {
	GraphDataDir string
	Arch         string
}

func NewCincinnati(log clog.PluggableLoggerInterface, config *v2alpha1.ImageSetConfiguration, opts mirror.CopyOptions, c Client, b bool, sig SignatureInterface) CincinnatiInterface {
	return &CincinnatiSchema{Log: log, Config: config, Opts: opts, Client: c, Fail: b, Signature: sig}
}

func (o *CincinnatiSchema) NewOCPClient() error {
	client, err := NewOCPClient(uuid.New(), o.Log)
	o.Client = client
	return err
}

func (o *CincinnatiSchema) NewOKDClient() error {
	client, err := NewOKDClient(uuid.New())
	o.Client = client
	return err
}

func (o *CincinnatiSchema) GetReleaseReferenceImages(ctx context.Context) ([]v2alpha1.CopyImageSchema, error) {
	cincinnatiParams := CincinnatiParams{
		GraphDataDir: filepath.Join(o.Opts.Global.WorkingDir, releaseImageExtractDir, cincinnatiGraphDataDir),
	}

	var (
		allImages  []v2alpha1.CopyImageSchema
		errs       = []error{}
		flagReport = false
	)

	// before making a deep copy
	// check that the "platform.release" field is not empty
	if len(o.Config.Mirror.Platform.Release) > 0 {
		copyImage := v2alpha1.CopyImageSchema{
			Source:      o.Config.Mirror.Platform.Release,
			Destination: "",
			Origin:      o.Config.Mirror.Platform.Release,
		}
		allImages = append(allImages, copyImage)
		return allImages, nil
	}

	filterCopy := o.Config.Mirror.Platform.DeepCopy()

	for _, arch := range filterCopy.Architectures {
		cincinnatiParams.Arch = arch
		o.CincinnatiParams = cincinnatiParams
		versionsByChannel := make(map[string]v2alpha1.ReleaseChannel, len(filterCopy.Channels))
		for _, ch := range filterCopy.Channels {
			var err error
			switch ch.Type {
			case v2alpha1.TypeOCP:
				err = o.NewOCPClient()
				if err != nil {
					errs = append(errs, err)
				}
			case v2alpha1.TypeOKD:
				err = o.NewOKDClient()
				if err != nil {
					errs = append(errs, err)
				}
			default:
				errs = append(errs, fmt.Errorf("invalid platform type %v", ch.Type))
				continue
			}
			if err != nil {
				errs = append(errs, err)
				continue
			}

			// CLID-135
			// detect and log as early as possible
			if len(ch.MaxVersion) > 0 && len(ch.MinVersion) > 0 {
				max := semver.MustParse(ch.MaxVersion)
				min := semver.MustParse(ch.MinVersion)
				if strings.Contains(ch.Name, "eus") && ((max.Minor - min.Minor) >= 2) && !flagReport {
					msg := "Extended Update Support (EUS) channel detected with minor version range >= 2\n" +
						"\t\t\t\tPlease refer to the web console https://access.redhat.com/labs/ocpupgradegraph/update_path\n" +
						"\t\t\t\tTo correctly determine the upgrade path for EUS releases"
					flagReport = true
					o.Log.Warn(msg)
				}
			}

			if len(ch.MaxVersion) == 0 || len(ch.MinVersion) == 0 {
				// Find channel maximum value and only set the minimum as well if heads-only is true
				if len(ch.MaxVersion) == 0 {
					latest, err := GetChannelMinOrMax(ctx, *o, ch.Name, false)
					if err != nil {
						errs = append(errs, err)
						continue
					}

					// Update version to release channel
					ch.MaxVersion = latest.String()
					o.Log.Debug("detected minimum version as %s", ch.MaxVersion)
					if len(ch.MinVersion) == 0 && ch.IsHeadsOnly() {
						//min, found := prevChannels[ch.Name]
						//if !found {
						// Starting at a new headsOnly channels
						min := latest.String()
						//}
						ch.MinVersion = min
						o.Log.Debug("detected minimum version as %s\n", ch.MinVersion)
					}
				}

				// Find channel minimum if full is true or just the minimum is not set
				// in the config
				if len(ch.MinVersion) == 0 {
					first, err := GetChannelMinOrMax(ctx, *o, ch.Name, true)
					if err != nil {
						errs = append(errs, err)
						continue
					}
					ch.MinVersion = first.String()
					o.Log.Debug("detected minimum version as %s\n", ch.MinVersion)
				}
				versionsByChannel[ch.Name] = ch
			} else {
				// Range is set. Ensure full is true so this
				// is skipped when processing release metadata.
				o.Log.Debug("processing minimum version %s and maximum version %s", ch.MinVersion, ch.MaxVersion)
				ch.Full = true
				versionsByChannel[ch.Name] = ch
			}

			downloads, err := getChannelDownloads(ctx, *o, nil, ch)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			allImages = append(allImages, downloads...)
		}

		// Update cfg release channels with maximum and minimum versions
		// if applicable
		for i, ch := range filterCopy.Channels {
			ch, found := versionsByChannel[ch.Name]
			if found {
				filterCopy.Channels[i] = ch
			}
		}

		if len(filterCopy.Channels) > 1 {
			newDownloads, err := getCrossChannelDownloads(ctx, *o, filterCopy.Channels)
			if err != nil {
				errs = append(errs, fmt.Errorf("[GetReleaseReferenceImages] error calculating cross channel upgrades: %v", err))
				continue
			}
			allImages = append(allImages, newDownloads...)
		}
	}

	imgs, err := o.Signature.GenerateReleaseSignatures(ctx, allImages)
	if err != nil {
		o.Log.Error("%v", err)
	}

	errorArray := []string{}
	for _, e := range errs {
		errorArray = append(errorArray, e.Error())
	}
	if len(errs) > 0 {
		return imgs, fmt.Errorf("[GetReleaseReferenceImages] error list %v", errorArray)
	}
	return imgs, nil
}

// getDownloads will prepare the downloads map for mirroring
func getChannelDownloads(ctx context.Context, cs CincinnatiSchema, lastChannels []v2alpha1.ReleaseChannel, channel v2alpha1.ReleaseChannel) ([]v2alpha1.CopyImageSchema, error) {
	var allImages []v2alpha1.CopyImageSchema

	var prevChannel v2alpha1.ReleaseChannel
	for _, ch := range lastChannels {
		if ch.Name == channel.Name {
			prevChannel = ch
		}
	}
	cs.Log.Trace("previous channel %v", prevChannel)
	// Plot between min and max of channel
	first, err := semver.Parse(channel.MinVersion)
	if err != nil {
		return allImages, err
	}
	last, err := semver.Parse(channel.MaxVersion)
	if err != nil {
		return allImages, err
	}

	var newDownloads []v2alpha1.CopyImageSchema
	if channel.ShortestPath {
		current, newest, updates, err := CalculateUpgrades(ctx, cs, channel.Name, channel.Name, first, last)
		if err != nil {
			return allImages, err
		}
		newDownloads = gatherUpdates(cs.Log, current, newest, updates)

	} else {
		lowRange, err := semver.ParseRange(fmt.Sprintf(">=%s", first))
		if err != nil {
			return allImages, err
		}
		highRange, err := semver.ParseRange(fmt.Sprintf("<=%s", last))
		if err != nil {
			return allImages, err
		}
		versions, err := GetUpdatesInRange(ctx, cs, channel.Name, highRange.AND(lowRange))
		if err != nil {
			return allImages, err
		}
		newDownloads = gatherUpdates(cs.Log, Update{}, Update{}, versions)
	}
	allImages = append(allImages, newDownloads...)

	return allImages, nil
}

// getCrossChannelDownloads will determine required downloads between channel versions (for OCP only)
func getCrossChannelDownloads(ctx context.Context, cs CincinnatiSchema, channels []v2alpha1.ReleaseChannel) ([]v2alpha1.CopyImageSchema, error) {
	// Strip any OKD channels from the list

	var ocpChannels []v2alpha1.ReleaseChannel
	for _, ch := range channels {
		if ch.Type == v2alpha1.TypeOCP {
			ocpChannels = append(ocpChannels, ch)
		}
	}
	// If no other channels exist, return no downloads
	if len(ocpChannels) == 0 {
		return []v2alpha1.CopyImageSchema{}, nil
	}

	firstCh, first, err := FindRelease(ocpChannels, true)
	if err != nil {
		return []v2alpha1.CopyImageSchema{}, fmt.Errorf("failed to find minimum release version: %v", err)
	}
	lastCh, last, err := FindRelease(ocpChannels, false)
	if err != nil {
		return []v2alpha1.CopyImageSchema{}, fmt.Errorf("failed to find maximum release version: %v", err)
	}
	current, newest, updates, err := CalculateUpgrades(ctx, cs, firstCh, lastCh, first, last)
	if err != nil {
		return []v2alpha1.CopyImageSchema{}, fmt.Errorf("failed to get upgrade graph: %v", err)
	}
	return gatherUpdates(cs.Log, current, newest, updates), nil
}

// gatherUpdates
func gatherUpdates(log clog.PluggableLoggerInterface, current, newest Update, updates []Update) []v2alpha1.CopyImageSchema {
	var allImages []v2alpha1.CopyImageSchema
	uniqueImages := make(map[v2alpha1.CopyImageSchema]bool)

	for _, update := range updates {
		log.Debug("Found update %s", update.Version)
		uniqueImages[v2alpha1.CopyImageSchema{Source: update.Image, Destination: ""}] = true
	}

	if current.Image != "" {
		uniqueImages[v2alpha1.CopyImageSchema{Source: current.Image, Destination: ""}] = true
	}

	if newest.Image != "" {
		uniqueImages[v2alpha1.CopyImageSchema{Source: newest.Image, Destination: ""}] = true
	}

	for img := range uniqueImages {
		allImages = append(allImages, img)
	}
	return allImages
}
