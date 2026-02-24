from typing import Dict, List, Union
import requests  # pyright: ignore[reportMissingModuleSource]
import re

# 读取 m3u8.log 文件内容
with open('../m3u8.log', 'r') as f:
    m3u8_content = f.read()

# 提取所有 .ts 文件的 URL 和时长
ts_entries = re.findall(r'#EXTINF:([\d.]+),\s*(https://[^\s]+\.ts)', m3u8_content)

# 构造请求数据
test_data: Dict[str, Union[str, List[Dict[str, Union[int, float, str]]]]] = {
    "video_id": "fake225",
    "key": "564b434876433962314a5056414c6665",
    "iv": "00000000000000000000000000000000",
    "ts_data": [
        {
            "ts_path": url,
            "ts_sequence": idx,
            "duration": float(duration)  # 添加时长字段
            # "key": "564b434876433962314a5056414c6665",
            # "iv": "00000000000000000000000000000000"
        }
        for idx, (duration, url) in enumerate(ts_entries)
        if float(duration) > 0  # 确保时长大于 0
    ]
}

print(test_data)

# 发送请求
response = requests.post(
    "http://127.0.0.1:8088/video_ts/save",
    json=test_data,
    headers={"Content-Type": "application/json"}
)

# 检查响应状态和内容
print(f"Status Code: {response.status_code}")
print(f"Response Text: {response.text}")

# 尝试解析 JSON，如果失败则打印错误信息
try:
    if response.text.strip():  # 确保响应不为空
        print(f"Response JSON: {response.json()}")
    else:
        print("Empty response received")
except requests.exceptions.JSONDecodeError as e:
    print(f"JSON decode error: {e}")
    print(f"Raw response content: {response.text}")
