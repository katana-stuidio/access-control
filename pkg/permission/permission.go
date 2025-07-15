package permission

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	openfga "github.com/openfga/go-sdk"
)

// PermissionService lida com autorização via OpenFGA
type PermissionService struct {
	FGAClient *openfga.APIClient
	StoreID   string
}

// NewPermissionService instancia um novo serviço de permissão
func NewPermissionService(ctx context.Context, apiURL, storeID, apiToken string) (*PermissionService, error) {
	// Criação do cliente OpenFGA
	cfg, err := openfga.NewConfiguration(openfga.Configuration{
		ApiUrl: apiURL,
	})
	if err != nil {
		return nil, err
	}
	// Se usar autenticação via token
	if apiToken != "" {
		cfg.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	}
	client := openfga.NewAPIClient(cfg)
	return &PermissionService{
		FGAClient: client,
		StoreID:   storeID,
	}, nil
}

// AddRelation adiciona uma relação entre user e tenant (ex: estudante, professor, etc)
func (ps *PermissionService) AddRelation(ctx context.Context, userID, tenantID, relation string) error {
	req := openfga.WriteRequest{
		Writes: &openfga.WriteRequestWrites{
			TupleKeys: []openfga.TupleKey{
				{
					User:     fmt.Sprintf("user:%s", userID),
					Relation: relation, // ex: "estudante"
					Object:   fmt.Sprintf("tenant:%s", tenantID),
				},
			},
		},
	}
	_, _, err := ps.FGAClient.OpenFgaApi.Write(ctx, ps.StoreID).Body(req).Execute()
	return err
}

// RemoveRelation remove uma relação entre user e tenant
func (ps *PermissionService) RemoveRelation(ctx context.Context, userID, tenantID, relation string) error {
	req := openfga.WriteRequest{
		Deletes: &openfga.WriteRequestDeletes{
			TupleKeys: []openfga.TupleKeyWithoutCondition{
				{
					User:     fmt.Sprintf("user:%s", userID),
					Relation: relation,
					Object:   fmt.Sprintf("tenant:%s", tenantID),
				},
			},
		},
	}
	_, _, err := ps.FGAClient.OpenFgaApi.Write(ctx, ps.StoreID).Body(req).Execute()
	return err
}

// CheckPermission verifica se o user tem permissão (relation) sobre o tenant
func (ps *PermissionService) CheckPermission(ctx context.Context, userID, tenantID, relation string) (bool, error) {
	req := openfga.CheckRequest{
		TupleKey: openfga.CheckRequestTupleKey{
			User:     fmt.Sprintf("user:%s", userID),
			Relation: relation, // ex: "can_responder_questionarios"
			Object:   fmt.Sprintf("tenant:%s", tenantID),
		},
	}
	resp, _, err := ps.FGAClient.OpenFgaApi.Check(ctx, ps.StoreID).Body(req).Execute()
	if err != nil {
		return false, err
	}
	return resp.GetAllowed(), nil
}

// RegisterModelFromJSON registra um modelo de autorização a partir de JSON
func (ps *PermissionService) RegisterModelFromJSON(ctx context.Context, modelJSON string) (string, error) {
	var modelRequest openfga.WriteAuthorizationModelRequest

	if err := json.Unmarshal([]byte(modelJSON), &modelRequest); err != nil {
		return "", fmt.Errorf("erro ao fazer parse do JSON: %v", err)
	}

	resp, _, err := ps.FGAClient.OpenFgaApi.WriteAuthorizationModel(ctx, ps.StoreID).Body(modelRequest).Execute()
	if err != nil {
		return "", fmt.Errorf("erro ao registrar modelo: %v", err)
	}

	return resp.GetAuthorizationModelId(), nil
}

// UpdateModelFromJSON atualiza um modelo de autorização existente a partir de JSON
func (ps *PermissionService) UpdateModelFromJSON(ctx context.Context, modelID, modelJSON string) error {
	var modelRequest openfga.WriteAuthorizationModelRequest

	if err := json.Unmarshal([]byte(modelJSON), &modelRequest); err != nil {
		return fmt.Errorf("erro ao fazer parse do JSON: %v", err)
	}

	_, _, err := ps.FGAClient.OpenFgaApi.WriteAuthorizationModel(ctx, ps.StoreID).Body(modelRequest).Execute()
	if err != nil {
		return fmt.Errorf("erro ao atualizar modelo: %v", err)
	}

	return nil
}

// GetDefaultEducationalModel retorna o modelo padrão para o sistema educacional
func GetDefaultEducationalModel() string {
	return `{
		"schema_version": "1.1",
		"type_definitions": [
			{
				"type": "user",
				"relations": {
					"can_create_quiz": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "professor"
									}
								},
								{
									"computedUserset": {
										"relation": "admin"
									}
								}
							]
						}
					},
					"can_read_quiz": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "professor"
									}
								},
								{
									"computedUserset": {
										"relation": "estudante"
									}
								},
								{
									"computedUserset": {
										"relation": "admin"
									}
								},
								{
									"computedUserset": {
										"relation": "instituicao"
									}
								},
								{
									"computedUserset": {
										"relation": "secretaria"
									}
								},
								{
									"computedUserset": {
										"relation": "grupo_educacional"
									}
								},
								{
									"computedUserset": {
										"relation": "coordenador"
									}
								}
							]
						}
					},
					"can_update_quiz": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "professor"
									}
								},
								{
									"computedUserset": {
										"relation": "admin"
									}
								}
							]
						}
					},
					"can_delete_quiz": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "professor"
									}
								},
								{
									"computedUserset": {
										"relation": "admin"
									}
								}
							]
						}
					},
					"can_respond_quiz": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "estudante"
									}
								}
							]
						}
					},
					"can_view_rankings": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "estudante"
									}
								},
								{
									"computedUserset": {
										"relation": "instituicao"
									}
								},
								{
									"computedUserset": {
										"relation": "admin"
									}
								},
								{
									"computedUserset": {
										"relation": "secretaria"
									}
								},
								{
									"computedUserset": {
										"relation": "grupo_educacional"
									}
								},
								{
									"computedUserset": {
										"relation": "coordenador"
									}
								}
							]
						}
					},
					"can_view_teachers": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "instituicao"
									}
								},
								{
									"computedUserset": {
										"relation": "admin"
									}
								},
								{
									"computedUserset": {
										"relation": "secretaria"
									}
								},
								{
									"computedUserset": {
										"relation": "grupo_educacional"
									}
								},
								{
									"computedUserset": {
										"relation": "coordenador"
									}
								}
							]
						}
					},
					"can_view_students": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "instituicao"
									}
								},
								{
									"computedUserset": {
										"relation": "admin"
									}
								},
								{
									"computedUserset": {
										"relation": "secretaria"
									}
								},
								{
									"computedUserset": {
										"relation": "grupo_educacional"
									}
								}
							]
						}
					},
					"can_edit_profile": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "estudante"
									}
								},
								{
									"computedUserset": {
										"relation": "professor"
									}
								},
								{
									"computedUserset": {
										"relation": "admin"
									}
								}
							]
						}
					},
					"can_view_quizzes_by_class": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "professor"
									}
								},
								{
									"computedUserset": {
										"relation": "coordenador"
									}
								},
								{
									"computedUserset": {
										"relation": "admin"
									}
								}
							]
						}
					},
					"can_view_quizzes_by_student": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "professor"
									}
								},
								{
									"computedUserset": {
										"relation": "coordenador"
									}
								},
								{
									"computedUserset": {
										"relation": "admin"
									}
								}
							]
						}
					},
					"can_view_graphs": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "instituicao"
									}
								},
								{
									"computedUserset": {
										"relation": "admin"
									}
								},
								{
									"computedUserset": {
										"relation": "secretaria"
									}
								},
								{
									"computedUserset": {
										"relation": "grupo_educacional"
									}
								}
							]
						}
					},
					"can_access_full_system": {
						"computedUserset": {
							"relation": "admin"
						}
					},
					"professor": {},
					"estudante": {},
					"instituicao": {},
					"admin": {},
					"secretaria": {},
					"grupo_educacional": {},
					"coordenador": {}
				}
			},
			{
				"type": "quiz",
				"relations": {
					"owner": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "professor"
									}
								},
								{
									"computedUserset": {
										"relation": "admin"
									}
								}
							]
						}
					},
					"can_read": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "owner"
									}
								},
								{
									"computedUserset": {
										"relation": "estudante"
									}
								},
								{
									"computedUserset": {
										"relation": "admin"
									}
								},
								{
									"computedUserset": {
										"relation": "instituicao"
									}
								},
								{
									"computedUserset": {
										"relation": "secretaria"
									}
								},
								{
									"computedUserset": {
										"relation": "grupo_educacional"
									}
								},
								{
									"computedUserset": {
										"relation": "coordenador"
									}
								}
							]
						}
					},
					"can_respond": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "estudante"
									}
								}
							]
						}
					}
				}
			},
			{
				"type": "tenant",
				"relations": {
					"professor": {},
					"estudante": {},
					"instituicao": {},
					"admin": {},
					"secretaria": {},
					"grupo_educacional": {},
					"coordenador": {}
				}
			}
		]
	}`
}

// RegisterDefaultEducationalModel registra o modelo padrão do sistema educacional
func (ps *PermissionService) RegisterDefaultEducationalModel(ctx context.Context) (string, error) {
	modelJSON := GetDefaultEducationalModel()
	return ps.RegisterModelFromJSON(ctx, modelJSON)
}

// Exemplo de inicialização do PermissionService
func ExampleInitPermissionService() (*PermissionService, error) {
	apiURL := os.Getenv("OPENFGA_API_URL")     // Ex: "http://localhost:8080"
	storeID := os.Getenv("OPENFGA_STORE_ID")   // Coloque o store_id criado no OpenFGA
	apiToken := os.Getenv("OPENFGA_API_TOKEN") // Se necessário
	ctx := context.Background()
	return NewPermissionService(ctx, apiURL, storeID, apiToken)
}
