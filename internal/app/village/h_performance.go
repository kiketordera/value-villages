package village

import (
	"github.com/gin-gonic/gin"
)

// Shows the performance of the Service chosen (or automatically choses the only one if there is only one)
func (VS *Server) showPerformanceLocalService(c *gin.Context) {
	service, b := VS.getOurService(c, "h_performance", "showPerformanceLocalService", "service, _ := VS.getOurService")
	if !b {
		return
	}
	sales := VS.getAllSalesFromSpecificService(c, service.ID, "h_performance", "showPerformanceLocalService", "sales := VS.getAllSalesFromSpecificService")
	payments := VS.getAllPaymentsFromSpecificService(c, service.ID, "h_performance", "showPerformanceLocalService", "payments := VS.getAllPaymentsFromSpecificService")
	purchases := VS.getAllPurchasesFromSpecificService(c, service.ID, "h_performance", "showPerformanceLocalService", "purchases := VS.getAllPurchasesFromSpecificService")
	increase, decrease, total := VS.getTotal(c, sales, purchases, payments)
	workers := VS.getWorkersFromSpecificService(c, service.ID, "h_performance", "showPerformanceLocalService", "workers := VS.getWorkersFromSpecificService")

	render(c, gin.H{
		"increase":    increase,
		"decrease":    decrease,
		"total":       total,
		"workers":     workers,
		"serviceType": service.TypeID,
	}, "performance-service.html")
}

// Shows the performance of all Services from the type given
func (VS *Server) showPerformanceAllTypeService(c *gin.Context) {
	service, b := VS.getOurService(c, "h_performance", "showPerformanceAllTypeService", "service, b := VS.getOurService")
	if !b {
		return
	}
	sales := VS.getAllSalesFromTheSystemByServiceType(c, service.TypeID, "h_performance", "showPerformanceAllTypeService", "sales := VS.getAllSalesFromTheSystemByServiceType")
	payments := VS.getAllPaymentsFromTheSystemByServiceType(c, service.TypeID, "h_performance", "showPerformanceAllTypeService", "payments := VS.getAllPaymentsFromTheSystemByServiceType")
	purchases := VS.getAllStockGeneratedFromTheSystemByServiceType(c, service.TypeID, "h_performance", "showPerformanceAllTypeService", "purchases := VS.getAllStockGeneratedFromTheSystemByServiceType")

	increase, decrease, total := VS.getTotal(c, sales, purchases, payments)
	allServices := VS.getAllSpecificServiceFromSystem(c, service.TypeID, "h_performance", "showPerformanceAllTypeService", "allServices := VS.getAllSpecificServiceFromSystem")
	allServicesHTML := VS.servicesToHTML(c, allServices)
	render(c, gin.H{
		"villages":    allServicesHTML,
		"increase":    increase,
		"decrease":    decrease,
		"total":       total,
		"serviceType": service.TypeID,
	}, "performance-all-services.html")
}

// Shows the performance of the specific Worker
func (VS *Server) showPerformanceWorker(c *gin.Context) {
	service, b := VS.getOurService(c, "h_performance", "showPerformanceWorker", "service, b := VS.getOurService")
	if !b {
		return
	}
	id, _ := VS.getIDParamFromURL(c, "id")
	user := VS.getUserByID(c, id, "h_performance", "showPerformanceWorker", "user := VS.getUserByID")

	sales := VS.getAllSalesFromSpecificWorker(c, id, "h_performance", "showPerformanceWorker", "sales := VS.getAllSalesFromSpecificWorker")
	payments := VS.getAllPaymentsFromSpecificWorker(c, id, "h_performance", "showPerformanceWorker", "payments := VS.getAllPaymentsFromSpecificWorker")
	purchases := VS.getAllServicePurchasesFromSpecificWorker(c, id, "h_performance", "showPerformanceWorker", "purchases := VS.getAllServicePurchasesFromSpecificWorker")
	increase, decrease, total := VS.getTotalWorker(c, sales, purchases, payments)
	render(c, gin.H{
		"increase":    increase,
		"decrease":    decrease,
		"total":       total,
		"user":        user,
		"serviceType": service.TypeID,
	}, "performance-worker.html")
}

// Shows the performance of the specific service
func (VS *Server) showPerformanceSpecificService(c *gin.Context) {
	serviceID, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	service := VS.getServiceByID(c, serviceID, "h_performance", "showPerformanceSpecificService", "service := VS.getServiceByID", true)
	sales := VS.getAllSalesFromSpecificService(c, serviceID, "h_performance", "showPerformanceSpecificService", "sales := VS.getAllSalesFromSpecificService")
	payments := VS.getAllPaymentsFromSpecificService(c, serviceID, "h_performance", "showPerformanceSpecificService", "payments := VS.getAllPaymentsFromSpecificService")
	purchases := VS.getAllPurchasesFromSpecificService(c, serviceID, "h_performance", "showPerformanceSpecificService", "purchases := VS.getAllPurchasesFromSpecificService")
	increase, decrease, total := VS.getTotal(c, sales, purchases, payments)
	workers := VS.getWorkersFromSpecificService(c, serviceID, "h_performance", "showPerformanceSpecificService", "workers := VS.getWorkersFromSpecificService")
	render(c, gin.H{
		"increase":    increase,
		"decrease":    decrease,
		"total":       total,
		"wk":          service,
		"workers":     workers,
		"serviceType": service.Type,
	}, "performance-service.html")
}

/*************************  WORKER  *************************/

// Shows the performance page with the self information to the worker
func (VS *Server) showPerformanceSelf(c *gin.Context) {
	service, b := VS.getOurService(c, "h_performance", "showPerformanceSelf", "service, b := VS.getOurService")
	if !b {
		return
	}
	user := VS.getUserFomCookie(c, "h_performance", "showPerformanceSelf", "user := VS.getUserFomCookie")
	sales := VS.getAllSalesFromSpecificWorker(c, user.ID, "h_performance", "showPerformanceSelf", "sales := VS.getAllSalesFromSpecificWorker")
	payments := VS.getAllPaymentsFromSpecificWorker(c, user.ID, "h_performance", "showPerformanceSelf", "payments := VS.getAllPaymentsFromSpecificWorker")
	purchases := VS.getAllServicePurchasesFromSpecificWorker(c, user.ID, "h_performance", "showPerformanceSelf", "purchases := VS.getAllServicePurchasesFromSpecificWorker")
	increase, decrease, total := VS.getTotal(c, sales, purchases, payments)
	render(c, gin.H{
		"increase":    increase,
		"decrease":    decrease,
		"total":       total,
		"user":        user,
		"serviceType": service.TypeID,
	}, "performance-worker.html")
}
