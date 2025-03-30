from src import resistere_config
from src import app


def main():
    config = resistere_config.read_config()
    app.run(config)


if __name__ == "__main__":
    main()
