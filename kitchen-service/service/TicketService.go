package service

import (
	"fmt"
	"kitchen-service/dto"
	"kitchen-service/model"
	"kitchen-service/repository"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type TicketService struct {
	MenuItemRepo   *repository.MenuItemRepository
	RestaurantRepo *repository.RestaurantRepository
	TicketRepo     *repository.TicketRepository
}

func (service *TicketService) Verify(restaurantID string, items dto.TicketLineItemsDTO) bool {
	if !service.RestaurantRepo.ExistsById(restaurantID) {
		fmt.Println("Restaurant does not exist!")
		return false
	}
	for _, item := range items.TicketLineItems {
		if !service.MenuItemRepo.ExistsByIdAndRestaurantID(item.MenuItemId, restaurantID) {
			fmt.Println("Menu item does not exist in the restaurant!")
			return false
		}
	}
	return true
}

func (service *TicketService) Create(restaurantId string, orderId string, items dto.TicketLineItemsDTO) bool {
	if !service.RestaurantRepo.ExistsById(restaurantId) {
		fmt.Println(("Restaurant does not exist!"))
		return false
	} else {
		fmt.Println(("Restaurant found"))
		orderUuid, _ := uuid.Parse(orderId)

		restaurantUuid, _ := uuid.Parse(restaurantId)
		var it []model.TicketLineItem
		ticket := model.Ticket{ID: orderUuid, RestaurantID: restaurantUuid, TicketState: model.PENDING, Items: it}
		for _, item := range items.TicketLineItems {
			menuItem := service.MenuItemRepo.FindById(item.MenuItemId)
			fmt.Println(item.MenuItemId)
			ticketLineItem := model.TicketLineItem{MenuItem: menuItem, Quantity: item.Quantity}
			ticket.AddItem(ticketLineItem)

		}
		fmt.Println(ticket.RestaurantID)
		service.TicketRepo.CreateTicket(&ticket)
		return true

	}
}

func (service *TicketService) Update(ticketId string, ticketState string) error {
	id, err := uuid.Parse(ticketId)
	if err != nil {
		print(err)
		return err
	}
	var ticketStatus model.TicketState
	switch ticketState {
	case "pending":
		ticketStatus = model.PENDING
	case "accepted":
		ticketStatus = model.ACCEPTED
	case "rejected":
		ticketStatus = model.REJECTED
	}
	url := fmt.Sprintf("http://%s:%s/%s/%s", os.Getenv("ORDER_SERVICE_DOMAIN"), os.Getenv("ORDER_SERVICE_PORT"), ticketId, ticketState)
	resp, err := http.Get(url)
	if err != nil {
		print(err)
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("order client error")
	}

	service.TicketRepo.UpdateTicket(id, ticketStatus)
	return nil
}
