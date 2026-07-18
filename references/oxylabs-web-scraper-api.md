# Web Scraper API

[**Web Scraper API**](https://oxylabs.io/products/scraper-api/web) is an **all-in-one web data collection solution** designed for extracting real-time data at scale from any public website. It covers every stage of web scraping, from crawling URLs and improving access success rates to data parsing and delivery to your preferred storage, so you don’t have to manage proxies, request management, or infrastructure.

The tool is built to meet enterprise security standards, including SOC 2 Type II compliance, and offers fast-adapting infrastructure that dynamically adjusts to target websites to ensure high success rates and reliable data extraction across search engines, e-commerce sites, travel platforms, and more.

## Getting started

**Create your API user credentials**: Sign up for a free trial or purchase the product in the [**Oxylabs dashboard**](https://dashboard.oxylabs.io/en/registration) to create your API user credentials (`USERNAME` and `PASSWORD`).

### Request samples

The API manages proxy rotation, request retries, and handles bot monitoring systems automatically as part of its integrated infrastructure, so a single request is enough to retrieve all structured data.

Below, you'll find sample cURL requests. For examples in other programming languages, please refer to the relevant sections: [**Amazon**](/api-targets/e-commerce/amazon.md), [**Google**](/api-targets/search-engines/google.md), [**Other Websites**](/api-targets/any-domain.md).

{% tabs %}
{% tab title="Amazon" %}

```bash
curl 'https://realtime.oxylabs.io/v1/queries' \
--user "USERNAME:PASSWORD" \
-H "Content-Type: application/json" \
-d '{
        "source": "amazon_product",
        "query": "B07FZ8S74R",
        "geo_location": "90210",
        "parse": true
    }'
```

{% endtab %}

{% tab title="Google" %}

```bash
curl 'https://realtime.oxylabs.io/v1/queries' \
--user 'USERNAME:PASSWORD' \
-H 'Content-Type: application/json' \
-d '{
        "source": "google_search",
        "query": "adidas",
        "geo_location": "California,United States",
        "parse": true
    }'
```

{% endtab %}

{% tab title="Other" %}

```bash
curl 'https://realtime.oxylabs.io/v1/queries' \
--user 'USERNAME:PASSWORD' \
-H 'Content-Type: application/json' \
-d '{
        "source": "universal",
        "url": "https://sandbox.oxylabs.io/"
    }'
```

{% endtab %}
{% endtabs %}

We use synchronous [**Realtime**](/products/web-scraper-api/integration-methods/realtime.md) integration method in our examples. If you would like to use [**Proxy Endpoint**](/products/web-scraper-api/integration-methods/proxy-endpoint.md) or asynchronous [**Push-Pull**](/products/web-scraper-api/integration-methods/push-pull.md) integration, refer to the [**integration methods**](/products/web-scraper-api/integration-methods.md) section.

{% file src="/files/q6L6LYAaqoK4WmPt29Ty" %}

{% file src="/files/YikpYmKx2dVWhzeNXOTN" %}

### Request parameter values

1. <mark style="background-color:green;">**`source`**</mark> - This parameter sets the scraper that will be used to process your request.
2. <mark style="background-color:green;">**`URL`**</mark> or <mark style="background-color:green;">**`query`**</mark> - Provide the `URL` or `query` for the type of page you want to scrape. Refer to the table below and the corresponding target sub-pages for detailed guidance on when to use each parameter.
3. Optionally, you can include additional parameters such as `geo_location`, `user_agent_type`, `parse` (find the list of our parsers [**here**](/products/web-scraper-api/features/result-processing-and-storage/dedicated-parsers.md)), `render` and more to customize your scraping request. Read more: [**Features**](/products/web-scraper-api/features.md).

&#x20;    \- mandatory parameter

### Scraping with URLs or parametrized inputs

Oxylabs support two general groups of inputs - URLs and parametrized inputs like queries, product or video IDs. [Generic targets](/api-targets/any-domain.md) which do not have a dedicated source can be scraped with `universal` source.

<table><thead><tr><th width="217">Target</th><th width="246">Source (Scraping URL)</th><th>Source (Using Query, Product or Video ID)</th></tr></thead><tbody><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/R2fZ2ApzQ6dZ8ZMeOpcN"><strong>Amazon</strong></a></td><td><code>amazon</code></td><td><p><code>amazon_product</code>,</p><p><code>amazon_search</code>,</p><p><code>amazon_pricing</code>,</p><p><code>amazon_sellers</code>,</p><p><code>amazon_bestsellers</code></p></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/ISyQMQli8hlrwctJhsAr"><strong>Google</strong></a></td><td><code>google</code></td><td><p><code>google_search</code>,</p><p><code>google_ads</code>,</p><p><code>google_ai_mode</code>,</p><p><code>google_lens</code>,</p><p><code>google_maps</code>,</p><p><code>google_travel_hotels</code>,</p><p><code>google_trends_explore</code>,</p><p><code>google_shopping_product</code>,</p><p><code>google_shopping_search</code></p></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/73eNAKzhETgsM48cMmtk"><strong>Bing</strong></a></td><td><code>bing</code></td><td><code>bing_search</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/aZztMmhf4QpVNVmP6x5V"><strong>YouTube</strong></a></td><td><code>universal</code></td><td><p><code>youtube_search</code>,</p><p><code>youtube_search_max</code>,</p><p><code>youtube_video_trainability</code>,</p><p><code>youtube_download</code>,</p><p><code>youtube_subtitles</code>,</p><p><code>youtube_metadata</code>,</p><p><code>youtube_channel</code>,</p><p><code>youtube_autocomplete</code></p></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/aDj8WgHNtNIgmhbbOyZO"><strong>ChatGPT</strong></a></td><td><code>universal</code></td><td><code>chatgpt</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/VXVoGSoH5MTbqZ1gBkrT"><strong>Perplexity</strong></a></td><td><code>universal</code></td><td><code>perplexity</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/s5cbCZDUL9mO9ECi3Yrq"><strong>Walmart</strong></a></td><td><code>walmart</code></td><td><p><code>walmart_search</code>,</p><p><code>walmart_product</code></p></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/ziI5eUiqAzWXg9Va5UnT"><strong>TikTok</strong></a></td><td><code>universal</code></td><td><p><code>tiktok_shop_search</code>,</p><p><code>tiktok_shop_product</code></p></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/RbiUEvUhBi35bvgE68DO"><strong>eBay</strong></a></td><td><code>ebay</code></td><td><p><code>ebay_search</code>,</p><p><code>ebay_product</code></p></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/JwnjVDP9hJESHIMLYue4"><strong>Etsy</strong></a></td><td><code>etsy</code></td><td><p><code>etsy_search</code>,</p><p><code>etsy_product</code></p></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/v2oGWmbKVrWEfXtqSr8T"><strong>Best Buy</strong></a></td><td><code>universal</code></td><td><p><code>bestbuy_search</code>,</p><p><code>bestbuy_product</code></p></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/S9dB9MLguv3bbtb3tuCM"><strong>Bed Bath &#x26; Beyond</strong></a></td><td><code>bedbathandbeyond</code></td><td><code>bedbathandbeyond_search</code>,<br><code>bedbathandbeyond_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/8ATBFgFjGbPhjFqFcPxo"><strong>Bodega Aurrerá</strong></a></td><td><code>bodegaaurrera</code></td><td><code>bodegaaurrera_search</code>,<br><code>bodegaaurrera_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/zqz78AjHVpRC0nEs4amy"><strong>Instacart</strong></a></td><td><code>instacart</code></td><td><code>instacart_search</code>,<br><code>instacart_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/wsvlWLTPtHZkVow7nMuh"><strong>Kroger</strong></a></td><td><code>kroger</code></td><td><p><code>kroger_search</code>,</p><p><code>kroger_product</code></p></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/W3koIlyvjuPtJ8lVVcDU"><strong>Lowe's</strong></a></td><td><code>lowes</code></td><td><p><code>lowes_search</code>,</p><p><code>lowes_product</code></p></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/yi0e1U3b2MJhm68gxoKj"><strong>Publix</strong></a></td><td><code>publix</code></td><td><code>publix_search</code>,<br><code>publix_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/sHUqw1rTKUAcvhsVoCWZ"><strong>Target</strong></a></td><td><code>target</code></td><td><p><code>target_search</code>,</p><p><code>target_product</code>,</p><p><code>target_category</code></p></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/AoKiv9r4OSvlO6lZrBxt"><strong>Grainger</strong></a></td><td><code>grainger</code></td><td><code>grainger_search</code>,<br><code>grainger_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/Sm4zwjm6TLGlhxB3CDrz"><strong>Costco</strong></a></td><td><code>costco</code></td><td><p><code>costco_search</code>,</p><p><code>costco_product</code></p></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/I55pfvcEQMvsw8HRDPje"><strong>Menards</strong></a></td><td><code>menards</code></td><td><code>menards_search</code>,<br><code>menards_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/CWKeXHxApSMyMoKRPs1N"><strong>Petco</strong></a></td><td><code>universal</code></td><td><code>petco_search</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/nfe2igYG2QQOs4Oe4pnv"><strong>Staples</strong></a></td><td><code>universal</code></td><td><code>staples_search</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/Bntk4oUw3xNrlQLgzLYm"><strong>Allegro</strong></a></td><td><code>universal</code></td><td><p><code>allegro_search</code>,</p><p><code>allegro_product</code></p></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/zoqwXL3LEhep6jFCONRP"><strong>Idealo</strong></a></td><td><code>universal</code></td><td><code>idealo_search</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/i2A88MhJRHPJ1n5r3mXI"><strong>MediaMarkt</strong></a></td><td><code>mediamarkt</code></td><td><code>mediamarkt_search</code>,<br><code>mediamarkt_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/wu806Go4c2aLwyBkuZ3l"><strong>Cdiscount</strong></a></td><td><code>cdiscount</code></td><td><code>cdiscount_search</code>,<br><code>cdiscount_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/XlcEKz4hezgvEDStqqut"><strong>Alibaba</strong></a></td><td><code>alibaba</code></td><td><code>alibaba_search</code>,<br><code>alibaba_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/WlpWnraruvMhIM9QcxUY"><strong>AliExpress</strong></a></td><td><code>aliexpress</code></td><td><code>aliexpress_search</code>,<br><code>aliexpress_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/pebcQY8J8bDjScNd76bn"><strong>IndiaMART</strong></a></td><td><code>indiamart</code></td><td><code>indiamart_search</code>,<br><code>indiamart_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/GffUmy19QDWBdYlM9NGr"><strong>Avnet</strong></a></td><td><code>universal</code></td><td><code>avnet_search</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/RTAKf6amL46IuFpXSYFM"><strong>Lazada</strong></a></td><td><code>lazada</code></td><td><code>lazada_search</code>,<br><code>lazada_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/jgON5w97o3g0r6Topqtl"><strong>Rakuten</strong></a></td><td><code>universal</code></td><td><code>rakuten_search</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/LAmGxXQiX5jgB8K4NEhC"><strong>Tokopedia</strong></a></td><td><code>universal</code></td><td><code>tokopedia_search</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/ivYqfplVIiIBgoQ5K4P6"><strong>Flipkart</strong></a></td><td><code>flipkart</code></td><td><code>flipkart_search</code>,<br><code>flipkart_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/0MSE7ZwqnLKsOhCjDOHL"><strong>MercadoLibre</strong></a></td><td><code>universal</code></td><td><code>mercadolibre_search</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/6ZIfwXZsGm1WSWQnbML5"><strong>Mercado Livre</strong></a></td><td><code>universal</code></td><td><code>mercadolivre_search</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/2RlTrgygUDwKJzauybMF"><strong>Magazine Luiza</strong></a></td><td><code>magazineluiza</code></td><td><code>magazineluiza_search</code>,<br><code>magazineluiza_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/PhhbAaT5yDnTAZx7rdcC"><strong>Falabella</strong></a></td><td><code>falabella</code></td><td><code>falabella_search</code>,<br><code>falabella_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/IxALcUCDXN2gvOTz1z30"><strong>Dcard</strong></a></td><td><code>universal</code></td><td><code>dcard_search</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/oOSpMnzNh1YVaTomRRcq"><strong>Airbnb</strong></a></td><td><code>airbnb</code></td><td><code>airbnb_product</code></td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/GbsYaRZalFMzPanQ3LNK"><strong>Zillow</strong></a></td><td><code>zillow</code></td><td>Using <code>query</code> parameter is not supported</td></tr><tr><td><a href="/spaces/BaJCXoqO1zdFnaMnCKrg/pages/Tq4ZdEkk3kFieXQw24v3"><strong>Other Websites</strong></a></td><td><code>universal</code></td><td>Using <code>query</code> parameter is not supported</td></tr></tbody></table>

{% hint style="info" %}
If you need help making your first request or optimizing your setup, our 24/7 expert support team is available via live chat.
{% endhint %}

## Testing via Scraper APIs Playground

Try [**Web Scraper API**](https://oxylabs.io/products/scraper-api/web) and [**OxyCopilot**](https://oxylabs.io/products/scraper-api/ai-web-scraper-copilot) in the [**Scraper APIs Playground**](https://dashboard.oxylabs.io/?route=/api-playground).

{% embed url="<https://www.youtube.com/watch?v=kDhknxrod6U>" %}

{% embed url="<https://youtu.be/9JoF8_5r5HY?si=61c3Zkx6FrH06PVa>" %}

## Testing via Postman

Get started with our API using Postman, a handy tool for making HTTP requests. Download our [**Web Scraper API Postman collection**](https://files.gitbook.com/v0/b/gitbook-x-prod.appspot.com/o/spaces%2FzrXw45naRpCZ0Ku9AjY1%2Fuploads%2FMeGA0TZQMcAFHoVhRSQi%2FWeb%20Scraper%20API.new_postman_collection.json?alt=media\&token=9f51d41b-6604-4eef-b6c1-5024cf52c5bf) and import it. This collection includes examples that demonstrate the functionality of the scraper. Customize the examples to your needs or start scraping right away.

For step-by-step instructions, watch our video tutorial below. If you're new to Postman, check out this short [**guide**](/integrations/web-scraper-api-integrations/postman.md).

{% embed url="<https://www.youtube.com/watch?v=WOD0mZnu-j0>" %}

{% hint style="info" %}
*All information herein is provided on an "as is" basis and for informational purposes only. We make no representation and disclaim all liability with respect to your use of any information contained on this page. Before engaging in scraping activities of any kind you should consult your legal advisors and carefully read the particular website's terms of service or receive a scraping license.*
{% endhint %}
