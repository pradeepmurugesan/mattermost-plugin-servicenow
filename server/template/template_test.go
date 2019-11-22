package template

import (
	"github.com/mattermost/mattermost-plugin-servicenow/server/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Render Template")
}

var _ = Describe("Render Template", func() {

	It("should throw an error in case the template is not found", func() {

		_, err := RenderTemplate("Some-non-existing-template", &models.Incident{})

		Expect(err.Error()).To(Equal("Couldn't find the template with name Some-non-existing-template"))

	})

	It("should render the incident template", func() {

		var incident = models.Incident{
			SysCreatedBy:     "some-user",
			ShortDescription: "description of the ticket",
			CreatedByID:      "some-id",
			Priority:         "1",
			Impact:           "5",
		}
		renderedTemplate, err := RenderTemplate("incident", incident)

		Expect(err).To(BeNil())
		Expect(renderedTemplate).NotTo(BeEmpty())
		Expect(renderedTemplate).To(Equal("\n### New incident created\n\nCreated By : some-user\nImpact : 5\nPriority : 1\nDescription: description of the ticket\n"))

	})
})
