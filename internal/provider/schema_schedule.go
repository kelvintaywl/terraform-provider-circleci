package circleci

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetScheduleSchema() map[string]*schema.Schema {
	// TODO: add validation?

	return map[string]*schema.Schema{
		"id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": &schema.Schema{
			Type:        schema.TypeString,
			Description: "Name of scheduled pipeline",
			Required:    true,
		},
		"description": &schema.Schema{
			Type:        schema.TypeString,
			Description: "Description of scheduled pipeline",
			Required:    true,
		},
		"project_slug": &schema.Schema{
			Type:        schema.TypeString,
			Description: "Project slug for this scheduled pipeline",
			Required:    true,
		},
		"actor": &schema.Schema{
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"loging": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"attribution_actor": &schema.Schema{
			Type:     schema.TypeString,
			Computed: false,
			ExactlyOneOf: []string{
				"system",
				"current",
			},
		},
		// "parameters"
		// "timeline"
	}
}
