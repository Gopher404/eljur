from selenium import webdriver
import config.config as config


def get_token():
    print("set options")
    opt = webdriver.ChromeOptions()
    opt.add_argument(f"user-data-dir={config.chrome_user_data_path}")

    print("open chrome")
    with webdriver.Chrome(options=opt) as driver:
        print("get vk.com")
        driver.get("https://vk.com/")
        print("page tittle:", driver.title)

        with open("../../vk_get_token/script.js") as f:
            script = f.read()

        print("execute js")
        token = driver.execute_script(script)
        print("token:", token)
        print("CLOSED")

    return token
