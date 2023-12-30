from locust import HttpUser, task

'''
Install

     sudo apt install pip
     pip3 install locust

To Run:

     locust -f ./scripts/locustfile.py --host=http://localhost:3080

Then open browser (default link):

     http://0.0.0.0:8089

'''
class QuickstartUser(HttpUser):

    @task(1)
    def index_page(self):
        self.client.get("/api/v1/author")
