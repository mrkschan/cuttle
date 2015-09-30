Setup
-----

Setup Python environment.

```
virtualenv --clear env
env/bin/pip install ShopifyAPI
```

Prepare API credentials.

```
touch env.sh
echo "export API_KEY='YOUR_API_KEY'" >> env.sh
echo "export PASSWORD='YOUR_PASSWORD'" >> env.sh
echo "export SHOPIFY_DOMAIN='YOUR_SHOPIFY_STORE'" >> env.sh
```


Run
---

```
HTTPS_PROXY=127.0.0.1:3128 ./load.rb
```


Expected
--------

No WARN print to stderr.
