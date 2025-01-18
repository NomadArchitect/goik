// Copyright 2025 Hans Jørgen Grimstad
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package robot

// Note:
//
//	All servos are not created equal.
//	Some servos will have a center position at 0 degrees, others
//	may have the center position at a different angle
type ServoAngles struct {
	Coxa  float64 `json:"Coxa"`
	Femur float64 `json:"Femur"`
	Tibia float64 `json:"Tibia"`
}

func NewServoAngles(Coxa float64, Femur float64, Tibia float64) ServoAngles {
	return ServoAngles{Coxa: Coxa, Femur: Femur, Tibia: Tibia}
}
