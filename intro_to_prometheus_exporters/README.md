This talk was delivered at Grafana & Friends Ottawa, May 2 2024.

To view the slides:

1. install [present](https://pkg.go.dev/golang.org/x/tools/cmd/present):
    ```console
    $ go install golang.org/x/tools/cmd/present@latest
    ```
2. make sure the submodule's up-to-date:
    ```console
    $ git submodule update --recursive
    ```
3. run the present tool (with speaker notes if you want):
    ```console
    $ present -notes
    ```
4. open http://127.0.0.1:3999/pres.slide
