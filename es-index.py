#!/usr/bin/env python

import elasticsearch
import datetime
import random

if __name__ == '__main__':

    es_endpoint = "https://search-sturgill-hqdvjjhotc7ynshsnwd6kivtwm.us-east-1.es.amazonaws.com"

    es = elasticsearch.Elasticsearch([es_endpoint])

    for i in range(100):
        doc = {
            'author': 'winez',
            'text': 'Test data',
            'seqNum': i,
            'value': i%5,
            'rand': random.randrange(0, 21, 2),
            'timestamp': datetime.datetime.now().utcnow(),
        }

        res = es.index(index="es-index", body=doc)

