package permission

import (
	"context"
	"fmt"
	"log"
)

// ExampleUsage demonstra como usar o PermissionService para registrar modelos FGA
func ExampleUsage() {
	ctx := context.Background()

	// Inicializar o serviço de permissão
	ps, err := NewPermissionService(ctx, "http://localhost:8080", "your-store-id", "your-api-token")
	if err != nil {
		log.Fatal("Erro ao inicializar PermissionService:", err)
	}

	// Exemplo 1: Registrar o modelo padrão do sistema educacional
	modelID, err := ps.RegisterDefaultEducationalModel(ctx)
	if err != nil {
		log.Fatal("Erro ao registrar modelo padrão:", err)
	}
	fmt.Printf("Modelo registrado com ID: %s\n", modelID)

	// Exemplo 2: Registrar um modelo customizado a partir de JSON
	customModelJSON := `{
		"schema_version": "1.1",
		"type_definitions": [
			{
				"type": "document",
				"relations": {
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
										"relation": "editor"
									}
								}
							]
						}
					},
					"can_write": {
						"union": {
							"child": [
								{
									"computedUserset": {
										"relation": "owner"
									}
								},
								{
									"computedUserset": {
										"relation": "editor"
									}
								}
							]
						}
					},
					"owner": {},
					"editor": {}
				}
			}
		]
	}`

	customModelID, err := ps.RegisterModelFromJSON(ctx, customModelJSON)
	if err != nil {
		log.Fatal("Erro ao registrar modelo customizado:", err)
	}
	fmt.Printf("Modelo customizado registrado com ID: %s\n", customModelID)

	// Exemplo 3: Adicionar relações de usuário
	err = ps.AddRelation(ctx, "user123", "tenant456", "professor")
	if err != nil {
		log.Fatal("Erro ao adicionar relação:", err)
	}

	// Exemplo 4: Verificar permissões
	hasPermission, err := ps.CheckPermission(ctx, "user123", "tenant456", "can_create_quiz")
	if err != nil {
		log.Fatal("Erro ao verificar permissão:", err)
	}
	fmt.Printf("Usuário pode criar quiz: %t\n", hasPermission)

	// Exemplo 5: Remover relação
	err = ps.RemoveRelation(ctx, "user123", "tenant456", "professor")
	if err != nil {
		log.Fatal("Erro ao remover relação:", err)
	}
}

// ExampleEducationalPermissions demonstra as permissões do sistema educacional
func ExampleEducationalPermissions() {
	ctx := context.Background()

	ps, err := NewPermissionService(ctx, "http://localhost:8080", "your-store-id", "your-api-token")
	if err != nil {
		log.Fatal("Erro ao inicializar PermissionService:", err)
	}

	// Registrar o modelo educacional
	_, err = ps.RegisterDefaultEducationalModel(ctx)
	if err != nil {
		log.Fatal("Erro ao registrar modelo educacional:", err)
	}

	// Exemplos de permissões por papel:

	// Professor
	fmt.Println("=== Permissões do Professor ===")
	ps.AddRelation(ctx, "prof1", "tenant1", "professor")

	permissions := []string{
		"can_create_quiz",
		"can_read_quiz",
		"can_update_quiz",
		"can_delete_quiz",
		"can_view_quizzes_by_class",
		"can_view_quizzes_by_student",
		"can_edit_profile",
	}

	for _, permission := range permissions {
		hasPermission, _ := ps.CheckPermission(ctx, "prof1", "tenant1", permission)
		fmt.Printf("Professor - %s: %t\n", permission, hasPermission)
	}

	// Estudante
	fmt.Println("\n=== Permissões do Estudante ===")
	ps.AddRelation(ctx, "est1", "tenant1", "estudante")

	studentPermissions := []string{
		"can_read_quiz",
		"can_respond_quiz",
		"can_view_rankings",
		"can_edit_profile",
	}

	for _, permission := range studentPermissions {
		hasPermission, _ := ps.CheckPermission(ctx, "est1", "tenant1", permission)
		fmt.Printf("Estudante - %s: %t\n", permission, hasPermission)
	}

	// Admin
	fmt.Println("\n=== Permissões do Admin ===")
	ps.AddRelation(ctx, "admin1", "tenant1", "admin")

	adminPermissions := []string{
		"can_create_quiz",
		"can_read_quiz",
		"can_update_quiz",
		"can_delete_quiz",
		"can_view_rankings",
		"can_view_teachers",
		"can_view_students",
		"can_view_graphs",
		"can_access_full_system",
	}

	for _, permission := range adminPermissions {
		hasPermission, _ := ps.CheckPermission(ctx, "admin1", "tenant1", permission)
		fmt.Printf("Admin - %s: %t\n", permission, hasPermission)
	}
}
