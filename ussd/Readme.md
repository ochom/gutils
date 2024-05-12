# USSD Walker

This package helps you to easily create `USSD` steps and define rules of how to walk each step

## Installation

This package request [`go 1.22`](https://go.dev/dl/) and above

To install the package run below command in your project root folder

```sh
go get github.com/ochom/gutils/ussd
```

Thats it 😎. Now you can create your menus and launch your app

## Creating a menu

```go
package main

import (
  "github.com/ochom/gutils/ussd"
  "github.com/gofiber/fiber/v2"
)

var mainMenu *ussd.Step

func init(){
  mainMenu = ussd.NewMenu(func(params map[string]string)string{
    // you can also add checks here and return a dynamic welcome message
    // e.g check if the dialer is a new user
    return "Hello, welcome to GUtils USSD.\n1. Say hello\n2. Say goodbye"
  })

  // add your steps to the main menu
  mainMenu.AddStep(ussd.Step{
    Key: "1",
    End: true,
    Menu: func(params map[string]string)string{
      // perform checks, get user and send personalized messages
      // ! this is a mock
      user := sql.FindOne(&models.User{PhoneNumber: params["phone_number"]})
      return fmt.Sprintf("Hello %s, How are you doing today", user.Name)
    }
  })

  mainMenu.AddStep(ussd.Step{
    Key: "2",
    End: true,
    Menu: func(params map[string]string)string{
      // perform checks, get user and send personalized messages
      // ! this is a mock
      user := sql.FindOne(&models.User{PhoneNumber: params["phone_number"]})
      return fmt.Sprintf("Goodbye %s", user.Name)
    }
  })
}

func main(){
  // launch the app, e.g using gofiber
  app := fiber.New()
  app.Post("/ussd", func(ctx *fiber.Ctx) error{
      var req map[string]string
      if err := ctx.BodyParser(&req); err != nil{
        return err
      }

      params := ussd.Params{
        SessionId: req["session_id"], // session id from your provider
        PhoneNumber: req["phone_number"], // dialer phone number
        Text: req["text"], // the dialed text i.e *1*2*5#
      }
    }
  })

  if err := app.Listen(":8080"); err != nil {
    logs.Fatal("server launch failed: %s", err.Error())
  }
}

```

## Explanation

The `mainStep` is just but a step of the `ussd`

A `Step` should have the following properties

| Property | Type     | Required | Description                                                                                                                  |
| -------- | -------- | -------- | ---------------------------------------------------------------------------------------------------------------------------- |
| Key      | string   | false    | The key of the step, when empty, the step is treated as a wild card and will match any input. For instance, the `mainStep`   |
| Menu     | function | true     | The function to call when the `Step` is matched during a walk. The function takes `map[string]string` and returns a `string` |
| End      | bool     | false    | Used to determine whether to end the walk at this step                                                                       |
| Children | []\*Step | false    | a list of children nested within a `Step`                                                                                    |

## What to expect in the params argues of a Menu func

| Item         | Type   | Description                                                            |
| ------------ | ------ | ---------------------------------------------------------------------- |
| session_id   | string | The session id                                                         |
| phone_number | string | session user phone number                                              |
| text         | string | The full ussd string i.e `*401*2*3#`                                   |
| input        | string | the last value the user input i.e in `*401*2*3#` the input will be `3` |

## Word

Feel free to use, contribute and improve this library at will.

Happy Coding :)
