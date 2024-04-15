package cart

import (
	"fmt"

	typePayload "github.com/FrancoRutigliano/ecom/types"
)

func getCartItemsIDs(items []typePayload.CartCheckoutItem) ([]int, error) {
	// create a slice of int's equal to the len of the items
	productsIds := make([]int, len(items))

	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity fot the product %d", item.ProductID)
		}
		productsIds[i] = item.ProductID
	}

	return productsIds, nil
}

func (h *Handler) createOrder(ps []typePayload.Product, items []typePayload.CartCheckoutItem, userID int) (int, float64, error) {
	productMap := make(map[int]typePayload.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}
	//check if the product's are in stock
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}
	//calculate the total price
	totalPrice := calculateTotalPrice(items, productMap)
	//reduce the total quantity of the product in our db
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		err := h.store.UpdateProduct(product)
		if err != nil {
			return 0, 0, err
		}
	}
	// create the order
	orderID, err := h.orderstore.CreateOrder(typePayload.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address",
	})
	if err != nil {
		return 0, 0, err
	}
	// create the order items
	for _, item := range items {
		h.orderstore.CreateOrderItem(typePayload.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []typePayload.CartCheckoutItem, products map[int]typePayload.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}
	// iterate on cartItems a proving if the product is available
	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductID)
		}

		//check the stock of each product compare to the item quantity requested
		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the quantity requested", product.Name)
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []typePayload.CartCheckoutItem, products map[int]typePayload.Product) float64 {
	var total float64

	for _, item := range cartItems {
		product := products[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}
	return total
}
