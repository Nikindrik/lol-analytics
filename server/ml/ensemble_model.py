import numpy as np
from feature_extractor import FeatureExtractor
from model_loader import ModelLoader


class EnsembleModel:
    def __init__(self, role_dir):
        loader = ModelLoader(role_dir)
        self.catboost = loader.load_catboost()
        self.xgboost = loader.load_xgboost()
        self.lightgbm = loader.load_lightgbm()
        self.scaler = loader.load_scaler()
        self.features = loader.load_features()
        self.extractor = FeatureExtractor(self.features)

    def _predict_catboost(self, X):
        if not self.catboost:
            return None, None
        proba = self.catboost.predict_proba(X)[0]
        return int(np.argmax(proba)), float(np.max(proba))

    def _predict_xgboost(self, X):
        if not self.xgboost:
            return None, None
        proba = self.xgboost.predict_proba(X)[0]
        return int(np.argmax(proba)), float(np.max(proba))

    def _predict_lightgbm(self, X):
        if not self.lightgbm:
            return None, None
        proba = self.lightgbm.predict(X)[0]
        return int(np.argmax(proba)), float(np.max(proba))

    def predict(self, data):
        X_raw = self.extractor.transform(data)
        X = self.scaler.transform(X_raw) if self.scaler else X_raw

        preds = []
        confs = []

        cb_p, cb_c = self._predict_catboost(X_raw)
        if cb_p is not None:
            preds.append(cb_p)
            confs.append(cb_c)

        xgb_p, xgb_c = self._predict_xgboost(X)
        if xgb_p is not None:
            preds.append(xgb_p)
            confs.append(xgb_c)

        lgb_p, lgb_c = self._predict_lightgbm(X_raw)
        if lgb_p is not None:
            preds.append(lgb_p)
            confs.append(lgb_c)

        if not preds:
            return -1, 0.0

        unique = list(set(preds))
        scores = dict.fromkeys(unique, 0)

        for p, c in zip(preds, confs, strict=False):
            scores[p] += c

        cluster = max(scores, key=scores.get)
        confidence = float(np.mean(confs))

        return cluster, confidence
