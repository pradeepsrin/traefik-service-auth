
**Sample config:**

    application-service-1:
      client_header_key: "X-Client-Id"
      secrets:
        service-client-1:
          - password1
          - password2
        service-client-2:
          - password3
          - password4 #could be from secret management


**Usage of a client**

    curl -X GET "http://application-service-1/test" -H "X-Client-Id: service-client-1" -H "Authorization: Basic password1"
