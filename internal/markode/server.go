package markode

import (
	"net/http"
)

func NewServer(words chan<- string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			pw := r.FormValue("password")
			if pw != "glitch" {
				http.Error(w, "Wrong password", http.StatusUnauthorized)
				return
			}
			words <- r.FormValue("corpus")
			http.Redirect(w, r, "/", http.StatusSeeOther)
		case http.MethodGet:
			w.Write([]byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Glitch</title>
			<style>
				* {
					font: inherit;
				}
				
				body {
					font-family: monospace;
					background-color: black;
					color: white;
				}

				form {
					display: flex;
					flex-direction: column;
					width: 100%;
					max-width: 400px;
					gap: 8px;
				}

				textarea, input[type=password] {
					border-radius: 4px;
					border: 1px solid currentColor;
					background-color: transparent;
					color: inherit;
					text-indent: 4px;
					padding: 4px 0;
				}

				button {
					border-radius: 8px;
				}
			</style>
		</head>
		<body>
			<form method="post" action="/">
				<input name="password" type="password" placeholder="Tape le mot de passe..." required />
				<textarea name="corpus" max="1024" required placeholder="Tape du texte..."></textarea>
				<button type="submit">Envoyer</button>
			</form>
		</body>
		</html>
			`))
		}
	})
}
