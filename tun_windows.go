//go:build windows
// +build windows

package main

import (
	"fmt"
	"net/netip"

	"yimiaoxiehou/sockstun/winipcfg"

	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wintun"
	"golang.zx2c4.com/wireguard/tun"
)

// createTUNDevice 创建一个TUN设备
func createTUNDevice(ifname string) (tun.Device, error) {
	// 定义一个唯一的GUID用于TUN设备
	id := &windows.GUID{
		Data1: 0x0000000,
		Data2: 0xFFFF,
		Data3: 0xFFFF,
		Data4: [8]byte{0xFF, 0xe9, 0x76, 0xe5, 0x8c, 0x74, 0x06, 0x3e},
	}
	// 卸载现有的Wintun驱动（如果存在）
	_ = wintun.Uninstall()
	// 创建TUN设备并返回
	return tun.CreateTUNWithRequestedGUID(ifname, id, 0)
}

// setIPAddress 为TUN设备设置IP地址
func setIPAddress(dev tun.Device, ipStr string) error {
	// 解析IP地址字符串
	ip, err := netip.ParsePrefix(ipStr)
	if err != nil {
		return fmt.Errorf("解析 IP 地址失败: %v", err)
	}

	// 获取原生TUN设备
	nativeTunDevice := dev.(*tun.NativeTun)
	// 获取设备的LUID（本地唯一标识符）
	link := winipcfg.LUID(nativeTunDevice.LUID())
	// 为设备设置IP地址
	return link.SetIPAddresses([]netip.Prefix{ip})
}
