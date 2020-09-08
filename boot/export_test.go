// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2014-2019 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package boot

import (
	"fmt"

	"github.com/snapcore/snapd/secboot"
	"github.com/snapcore/snapd/snap"
)

func NewCoreBootParticipant(s snap.PlaceInfo, t snap.Type, dev Device) *coreBootParticipant {
	bs, err := bootStateFor(t, dev)
	if err != nil {
		panic(err)
	}
	return &coreBootParticipant{s: s, bs: bs}
}

func NewCoreKernel(s snap.PlaceInfo, d Device) *coreKernel {
	return &coreKernel{s, bootloaderOptionsForDeviceKernel(d)}
}

type Trivial = trivial

func (m *Modeenv) WasRead() bool {
	return m.read
}

func (m *Modeenv) DeepEqual(m2 *Modeenv) bool {
	return m.deepEqual(m2)
}

var (
	MarshalModeenvEntryTo        = marshalModeenvEntryTo
	UnmarshalModeenvValueFromCfg = unmarshalModeenvValueFromCfg

	NewTrustedAssetsCache = newTrustedAssetsCache

	ObserveSuccessfulBootWithAssets = observeSuccessfulBootAssets
	SealKeyToModeenv                = sealKeyToModeenv
)

type BootAssetsMap = bootAssetsMap
type TrackedAsset = trackedAsset

func (t *TrackedAsset) Equals(blName, name, hash string) error {
	equal := t.hash == hash &&
		t.name == name &&
		t.blName == blName
	if !equal {
		return fmt.Errorf("not equal to bootloader %q tracked asset %v:%v", t.blName, t.name, t.hash)
	}
	return nil
}

func (o *TrustedAssetsInstallObserver) CurrentTrustedBootAssetsMap() BootAssetsMap {
	return o.currentTrustedBootAssetsMap()
}

func (o *TrustedAssetsInstallObserver) CurrentTrustedRecoveryBootAssetsMap() BootAssetsMap {
	return o.currentTrustedRecoveryBootAssetsMap()
}

func MockSecbootSealKey(f func(key secboot.EncryptionKey, params *secboot.SealKeyParams) error) (restore func()) {
	old := secbootSealKey
	secbootSealKey = f
	return func() {
		secbootSealKey = old
	}
}

func (o *TrustedAssetsUpdateObserver) InjectChangedAsset(blName, assetName, hash string, recovery bool) {
	ta := &trackedAsset{
		blName: blName,
		name:   assetName,
		hash:   hash,
	}
	if !recovery {
		o.changedAssets = append(o.changedAssets, ta)
	} else {
		o.seedChangedAssets = append(o.seedChangedAssets, ta)
	}
}
