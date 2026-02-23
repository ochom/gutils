// Package ussd provides a framework for building USSD (Unstructured Supplementary Service Data)
// menu applications. It supports nested menus, dynamic content, and session management.
//
// USSD is commonly used in mobile telecommunications for banking, airtime purchase,
// and other interactive services.
//
// Example usage:
//
//	// Define the main menu
//	mainMenu := ussd.NewStep("Welcome!\n1. Check Balance\n2. Transfer Money\n3. Exit")
//
//	// Add sub-steps
//	mainMenu.AddStep("1", ussd.NewStep(func(params map[string]string) string {
//		balance := getBalance(params["phone_number"])
//		return fmt.Sprintf("Your balance is: KES %s", balance)
//	}).End = true)
//
//	mainMenu.AddStep("2", transferMenu)
//	mainMenu.AddStep("3", ussd.NewStep("Thank you for using our service").End = true)
//
//	// Register the menu
//	ussd.New(mainMenu)
//
//	// Process incoming USSD request
//	http.HandleFunc("/ussd", func(w http.ResponseWriter, r *http.Request) {
//		params := ussd.Params{
//			SessionId:   r.FormValue("sessionId"),
//			PhoneNumber: r.FormValue("phoneNumber"),
//			Text:        r.FormValue("text"),
//		}
//
//		step, err := ussd.Parse(params)
//		if err != nil {
//			fmt.Fprint(w, "END Error occurred")
//			return
//		}
//
//		fmt.Fprint(w, step.GetResponse())
//	})
package ussd

// Params holds the incoming USSD request parameters.
type Params struct {
	// Text is the full USSD input string (e.g., "1*2*3" for nested selections)
	Text string
	// SessionId uniquely identifies the USSD session
	SessionId string
	// PhoneNumber is the subscriber's phone number in international format
	PhoneNumber string
}
