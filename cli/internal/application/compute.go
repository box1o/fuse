package application

import (
	"context"
	"fmt"
	"os"
	"strings"

	"fuse/cli/internal/ports"
	"fuse/pkg/computeapi"

	"github.com/google/uuid"
)

const AgentVersion = "0.1.0"

type ComputeService struct {
	gateway       ports.NodeGateway
	credentials   ports.CredentialStore
	installations ports.InstallationStore
	detector      ports.CapabilityDetector
}

func NewComputeService(gateway ports.NodeGateway, credentials ports.CredentialStore, installations ports.InstallationStore, detector ports.CapabilityDetector) *ComputeService {
	return &ComputeService{
		gateway:       gateway,
		credentials:   credentials,
		installations: installations,
		detector:      detector,
	}
}

func (s *ComputeService) Inspect(ctx context.Context) (computeapi.Capabilities, error) {
	return s.detector.Detect(ctx)
}

func (s *ComputeService) Register(ctx context.Context, name string, capabilities computeapi.Capabilities) (computeapi.Node, bool, error) {
	token, err := s.token()
	if err != nil {
		return computeapi.Node{}, false, err
	}
	installationID, err := s.installations.EnsureInstallationID()
	if err != nil {
		return computeapi.Node{}, false, fmt.Errorf("create installation identity: %w", err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		return computeapi.Node{}, false, err
	}
	if strings.TrimSpace(name) == "" {
		name = hostname
	}
	return s.gateway.Register(ctx, token, computeapi.RegisterNodeRequest{
		InstallationID: installationID,
		Name:           name,
		Hostname:       hostname,
		AgentVersion:   AgentVersion,
		Capabilities:   capabilities,
	})
}

func (s *ComputeService) List(ctx context.Context) ([]computeapi.Node, error) {
	token, err := s.token()
	if err != nil {
		return nil, err
	}
	return s.gateway.List(ctx, token)
}

func (s *ComputeService) Get(ctx context.Context, nodeID string) (computeapi.Node, error) {
	if _, err := uuid.Parse(nodeID); err != nil {
		return computeapi.Node{}, fmt.Errorf("invalid node ID")
	}
	token, err := s.token()
	if err != nil {
		return computeapi.Node{}, err
	}
	return s.gateway.Get(ctx, token, nodeID)
}

func (s *ComputeService) Update(ctx context.Context, nodeID string, request computeapi.UpdateNodeRequest) (computeapi.Node, error) {
	if _, err := uuid.Parse(nodeID); err != nil {
		return computeapi.Node{}, fmt.Errorf("invalid node ID")
	}
	if request.Name == nil && request.Disabled == nil {
		return computeapi.Node{}, fmt.Errorf("no update was provided")
	}
	token, err := s.token()
	if err != nil {
		return computeapi.Node{}, err
	}
	return s.gateway.Update(ctx, token, nodeID, request)
}

func (s *ComputeService) Delete(ctx context.Context, nodeID string) error {
	if _, err := uuid.Parse(nodeID); err != nil {
		return fmt.Errorf("invalid node ID")
	}
	token, err := s.token()
	if err != nil {
		return err
	}
	return s.gateway.Delete(ctx, token, nodeID)
}

func (s *ComputeService) token() (string, error) {
	token, err := s.credentials.Load()
	if err != nil {
		return "", fmt.Errorf("run 'fuse auth login' first")
	}
	return token, nil
}
