package sbi

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/free5gc/amf/internal/rms"
)

func (s *Server) getRMSRoutes() []Route {
	return []Route{
		{
			Name:    "root",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.String(http.StatusOK, "Hello World!")
			},
		},
		// add more Route based on provided spec
		// Get operation
		// Return the currently existing UE-RM subscriptions
		{
			Name:    "getSubscription", 
			Method:  http.MethodGet,
			Pattern: "/subscriptions/",
			APIFunc: func(c *gin.Context) {
				subs := s.rms.QueryAll()
				c.JSON(http.StatusOK, gin.H{
					"subscriptions": subs,
				})
			},
		},
		// Post operation
		// Create a new UE-RM subscription
		// The SubId in the POST request can be filled arbitrarily (e.g., sub-001). Your API handler is then responsible for populating the assigned SubId and returning it.
		{
			Name:    "postSubscription",
			Method:  http.MethodPost,
			Pattern: "/subscriptions/",
			APIFunc: func(c *gin.Context) {
				var sub rms.Subscription
				if err := c.ShouldBindJSON(&sub); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				created := s.rms.Add(sub)
				c.JSON(http.StatusCreated, created)
			},
		},
		// PUT Operation
		// Create or overwrite a UE-RM subscription
		// The relevant schemas are the same as those for the POST operation.
		{
			Name:    "putSubscription",
			Method:  http.MethodPut,
			Pattern: "/subscriptions/:subId/",
			APIFunc: func(c *gin.Context) {
				subId := c.Param("subId")
				var sub rms.Subscription
				if err := c.ShouldBindJSON(&sub); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				if s.rms.Modify(subId, sub) {
					c.JSON(http.StatusOK, gin.H{"message": "Subscription modified"})
				} else {
					s.rms.Add(sub)
					c.JSON(http.StatusCreated, gin.H{"message": "Subscription created"})
				}
			},
		},
		// Delete operation
		// Delete an existing UE-RM subscription.
		{
			Name:    "deleteSubscription",
			Method:  http.MethodDelete,
			Pattern: "/subscriptions/:subId/",
			APIFunc: func(c *gin.Context) {
				subId := c.Param("subId")
				if s.rms.Delete(subId) {
					// c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted"})
					c.Status(http.StatusNoContent)
				} else {
					c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
				}
			},
		},
		
	}
}
