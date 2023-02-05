## List of commands for generating new swagger 

- `swag init -g app.go -d ./actions --parseDependency true`


### Example of an documented function

```go
// List godoc
// GET
// @Summary      List gets all Characters. This function is mapped to the path
// @Description  Returns a JSON list of all the characters registered in the databases with the possibility of pagination
// @Tags         characters
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  models.Characters
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /characters [get]
func (v CharactersResource) List(c buffalo.Context) error {
    //...
}
```