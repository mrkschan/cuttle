import functools
import json
import os

import shopify


API_KEY = os.environ.get('API_KEY', '')
PASSWORD = os.environ.get('PASSWORD', '')
SHOPIFY_DOMAIN = os.environ.get('SHOPIFY_DOMAIN', '')


def patch_ssl():
    # Assign SSL context to pyactiveresource by monkey patching.
    from six.moves import urllib
    import ssl

    # Set custom trusted CA.
    context = ssl.create_default_context()
    context.load_verify_locations(cafile='cacert.pem')

    urllib.request.urlopen = functools.partial(urllib.request.urlopen,
                                               context=context)

try:
    # Patch SSL config in Python 2.7.9+.
    patch_ssl()
except:
    # SSLContext not available, do not verify server cert instead.
    pass


# Use Shopify private API.
shop_url = 'https://{}:{}@{}/admin'.format(API_KEY, PASSWORD, SHOPIFY_DOMAIN)
shopify.ShopifyResource.set_site(shop_url)

try:
    print json.dumps(shopify.Shop.current().to_dict())
except Exception as exc:
    print exc
