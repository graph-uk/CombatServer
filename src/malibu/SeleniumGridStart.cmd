start java -jar D:\Environment\Utils\Selenium\selenium-server-standalone-2.45.0.jar -role hub -multiWindow -browserSessionReuse
timeout 1
start java -jar D:\Environment\Utils\Selenium\selenium-server-standalone-2.45.0.jar -Dwebdriver.chrome.driver="D:\Environment\Utils\ChromeWebDriver\win32_2.14\chromedriver.exe" -role node -hub http://localhost:4444/grid/register -port 5555