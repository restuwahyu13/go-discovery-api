# Golang Service Discovery

This is basic concept for `service registry` and `discovery` if you want more complex use case , you can check this [video tutorial](https://www.youtube.com/watch?v=s3I1kKKfjtQ) or you can buy this [e-book](https://www.amazon.com/Consul-Up-Running-Luke-Kysow-ebook/dp/B0B2ZHXLCV) from amazon.

## Benefit Using Consul

- [x] Monitoring service is health or unhealth 
- [x] Auto failover with cross server when one of service is dead, you can using `pairing feature`
- [x] Network security communication using tls encryption
- [x] Discovery & Registry service, you can find service using id or name
- [x] Management access control using `acl feature`, like A allow connect to B and C deny connecto to A 
- [x] You can use Load Balancing, Rate Limit, Circuit Bracker etc
- [x] Cross cloud platform support etc and container like `docker` or `kubernetes` for production and complex use case , use `kubernetes` recommended
- [x] Etc 

## Simple Sample Policy

```hcl
service_prefix "prod-service-" {
  policy = "read"
}

node_prefix "" {
  policy = "read"
}
```

## Full Sample Policy

```hcl
acl      = "write"
operator = "write"
mesh     = "write"
peering  = "write"
keyring  = "write"

agent_prefix "" {
  policy = "write"
}

node_prefix "" {
  policy = "write"
}

identity_prefix "" {
  policy = "write"
}

key_prefix "" {
  policy = "write"
}

service_prefix "" {
  policy     = "write"
  intentions = "write"
}

event_prefix "" {
  policy = "write"
}

session_prefix "" {
  policy = "write"
}

query_prefix "" {
  policy = "write"
}
```

## Example Use Service Discovery Direct Connect To Consul


```js
import ConsulHashicorp from 'consul'
import axios from 'axios'

class Consul {
  static consul = new ConsulHashicorp({ host: 'localhost', port: '8500', secure: false })

  static getKV(key) {
    return Consul.consul.kv.get(key)
  }

  static async health(serviceName, options) {
    try {
      const list = await Consul.consul.agent.check.list(options)
      if (!list[serviceName]) {
        return
      }
      return list[serviceName]
    } catch (e) {
      throw new Error(e)
    }
  }

  static async service(serviceName, options) {
    try {
      const list = await Consul.consul.agent.service.list(options)
      if (!list[serviceName]) {
        return
      }
      return list[serviceName]
    } catch (e) {
      throw new Error(e)
    }
  }
}

;(async () => {
  const secretToken = '4cf4f739-e472-dfa5-280e-1c507d5ee326' // I use production secret token

  const svcGreatday = await Consul.service('prod-service-greatday', { token: secretToken })
  const svcMekari = await Consul.service('prod-service-mekari', { token: secretToken })

  axios
    .get(`http://${svcGreatday?.Address}:${svcGreatday?.Port}`)
    .then(({ data }) => console.log(`\n${svcGreatday.Service} Data:\n`, data))
    .catch((e) => console.error(`\n${svcGreatday?.Service} Error:\n`, e?.message || e.cause?.code))

  axios
    .get(`http://${svcMekari?.Address}:${svcMekari?.Port}`)
    .then(({ data }) => console.log(`\n${svcMekari?.Service} Data:\n`, data))
    .catch((e) => console.error(`\n${svcMekari?.Service} Error:\n`, e?.message || e.cause?.code))
})()

```
