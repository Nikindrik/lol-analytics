import numpy as np
from regression_model_loader import RegressionModelLoader

class RegressionEnsemble:
    def __init__(self, model_dir):
        loader = RegressionModelLoader(model_dir)
        self.catboost = loader.load_catboost_regressor()
        self.xgboost = loader.load_xgboost_regressor()
        self.lightgbm = loader.load_lightgbm_regressor()
        self.scaler = loader.load_scaler()
        self.features = loader.load_features()
        print(f"Loaded regression model with features: {self.features}")
    
    def _predict_catboost(self, X):
        if not self.catboost:
            return None
        pred = self.catboost.predict(X)[0]
        return float(pred)
    
    def _predict_xgboost(self, X):
        if not self.xgboost:
            return None
        pred = self.xgboost.predict(X)[0]
        return float(pred)
    
    def _predict_lightgbm(self, X):
        if not self.lightgbm:
            return None
        pred = self.lightgbm.predict(X)[0]
        return float(pred)
    
    def predict(self, X_raw):
        # Масштабируем если есть scaler
        X = self.scaler.transform(X_raw) if self.scaler else X_raw
        
        preds = []
        
        cb_pred = self._predict_catboost(X_raw)
        if cb_pred is not None:
            preds.append(cb_pred)
        
        xgb_pred = self._predict_xgboost(X)
        if xgb_pred is not None:
            preds.append(xgb_pred)
        
        lgb_pred = self._predict_lightgbm(X_raw)
        if lgb_pred is not None:
            preds.append(lgb_pred)
        
        if not preds:
            return 0.0
        
        return float(np.mean(preds))