terraform {
  required_providers {
    theoffice = {
      source = "anGie44/theoffice"
    }
  }
}

provider "theoffice" {
	endpoint = "https://theofficeapi-angelinepinilla.b4a.run"
}

data "theoffice_quotes" "example" {
	season = 1
	episode = 2
}

output "quote" {
	value = data.theoffice_quotes.example.quotes[0]
}
