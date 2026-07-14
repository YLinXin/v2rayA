package serverObj

import (
	"encoding/json"
	"testing"
)

func TestXHTTPExtra(t *testing.T) {
	link := `vless://qwe@qwe.tencentapp.cn:443?encryption=qwe.native.0rtt.100-111-1111.75-0-111.50-0-3333.44h-qwe&flow=xtls-rprx-vision-udp443&security=tls&sni=qwe.m3u8.qwe.pay.uruu.cn&fp=firefox&alpn=h3&insecure=0&allowInsecure=0&type=xhttp&path=/zones&mode=auto&extra={"xPaddingBytes":"100-1000","noGRPCHeader":false,"noSSEHeader":false,"scMaxEachPostBytes":"1000000","scMinPostsIntervalMs":"35","scMaxBufferedPosts":30,"scStreamUpServerSecs":"20-80","xmux":{"maxConcurrency":"8-16","maxConnections":"0","cMaxReuseTimes":"0","hMaxRequestTimes":"600-900","hMaxReusableSecs":"1800-3000","hKeepAlivePeriod":0}}#剩余流量：96.95 GB`

	v, err := ParseVlessURL(link)
	if err != nil {
		t.Fatalf("ParseVlessURL failed: %v", err)
	}

	if v.Net != "xhttp" {
		t.Errorf("expected Net to be xhttp, got %v", v.Net)
	}

	if v.XHTTPMode != "auto" {
		t.Errorf("expected XHTTPMode to be auto, got %v", v.XHTTPMode)
	}

	if len(v.XHTTPExtra) == 0 {
		t.Fatalf("expected XHTTPExtra to be parsed, but got empty")
	}

	// Verify XHTTPExtra is valid JSON
	var extra map[string]interface{}
	if err := json.Unmarshal(v.XHTTPExtra, &extra); err != nil {
		t.Fatalf("failed to unmarshal XHTTPExtra: %v", err)
	}

	if extra["xPaddingBytes"] != "100-1000" {
		t.Errorf("expected xPaddingBytes to be 100-1000, got %v", extra["xPaddingBytes"])
	}

	// Test ExportToURL
	exportedLink := v.ExportToURL()
	v2, err := ParseVlessURL(exportedLink)
	if err != nil {
		t.Fatalf("failed to parse exported link: %v", err)
	}

	if string(v2.XHTTPExtra) != string(v.XHTTPExtra) {
		t.Errorf("exported extra doesn't match original:\nGot: %s\nExpected: %s", string(v2.XHTTPExtra), string(v.XHTTPExtra))
	}

	// Test Configuration construction
	info := PriorInfo{
		Tag: "proxy",
	}
	config, err := v.Configuration(info)
	if err != nil {
		t.Fatalf("v.Configuration failed: %v", err)
	}

	xhttpSettings := config.CoreOutbound.StreamSettings.XHTTPSettings
	if xhttpSettings == nil {
		t.Fatalf("expected XHTTPSettings to be populated")
	}

	if string(xhttpSettings.Extra) != string(v.XHTTPExtra) {
		t.Errorf("XHTTPSettings.Extra doesn't match original:\nGot: %s\nExpected: %s", string(xhttpSettings.Extra), string(v.XHTTPExtra))
	}

	t.Logf("Success! Configured xhttp settings: %+v", xhttpSettings)
}
