import numpy as np


class FeatureExtractor:
    def __init__(self, features):
        self.features = features

    def transform(self, data):
        return np.array([[
            float(data.get(f, 0) or 0)
            for f in self.features
        ]])
