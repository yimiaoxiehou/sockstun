//go:build linux
// +build linux

package main

import (
	"fmt"
	"os/exec"

	"golang.zx2c4.com/wireguard/tun"
)

// createTUNDevice 创建一个TUN设备
func createTUNDevice(ifname string) (tun.Device, error) {
	// 使用wireguard-go库创建TUN设备
	return tun.CreateTUN(ifname, 0)
}

// setIPAddress 为TUN设备设置IP地址
func setIPAddress(dev tun.Device, ipStr string) error {
	// 获取TUN设备的名称
	name, err := dev.Name()
	if err != nil {
		return fmt.Errorf("获取设备名称失败: %v", err)
	}

	// 使用ip命令为设备设置IP地址
	cmd := exec.Command("ip", "addr", "add", ipStr, "dev", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("设置 IP 地址失败: %v", err)
	}

	return nil
}
