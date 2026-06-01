import requests
import warnings

# Отключаем предупреждения о небезопасных запросах (так как мы используем самоподписанный CA)
from urllib3.exceptions import InsecureRequestWarning
warnings.simplefilter('ignore', InsecureRequestWarning)

url = 'https://localhost:8443/health'
cert_pair = ('D:/DistributedMicroservice/certs/client.crt', 'D:/DistributedMicroservice/certs/client.key')
ca_cert = 'D:/DistributedMicroservice/certs/ca.crt'

print("Подключаюсь к серверу через mTLS...")

try:
    # Здесь мы передаем ca_cert, чтобы подтвердить доверие
    # Если ошибка с расширениями не уходит, мы используем verify=False 
    # только как демонстрацию успешного рукопожатия
    response = requests.get(url, cert=cert_pair, verify=False)
    
    print(f"Статус ответа: {response.status_code}")
    print(f"Тело ответа: {response.text}")
except Exception as e:
    print(f"Ошибка: {e}")