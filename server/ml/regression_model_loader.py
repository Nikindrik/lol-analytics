import json
import catboost
import joblib
import lightgbm as lgb
import xgboost as xgb
from pathlib import Path

class RegressionModelLoader:
    def __init__(self, model_dir):
        self.model_dir = Path(model_dir)
    
    def load_catboost_regressor(self):
        path = self.model_dir / 'catboost.cbm'
        if not path.exists():
            return None
        model = catboost.CatBoostRegressor()
        model.load_model(str(path), format='cbm')
        return model
    
    def load_xgboost_regressor(self):
        path = self.model_dir / 'xgboost.json'
        if not path.exists():
            return None
        model = xgb.XGBRegressor()
        model.load_model(str(path))
        return model
    
    def load_lightgbm_regressor(self):
        path = self.model_dir / 'lightgbm.txt'
        if not path.exists():
            return None
        return lgb.Booster(model_file=str(path))
    
    def load_scaler(self):
        path = self.model_dir / 'scaler.pkl'
        if not path.exists():
            return None
        return joblib.load(path)
    
    def load_features(self):
        path = self.model_dir / 'features.json'
        if not path.exists():
            return []
        with open(path) as f:
            return json.load(f)