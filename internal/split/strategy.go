package split

import (
	i3 "go.i3wm.org/i3/v4"
)

type SplitStrategy interface {
	SetUpNextSplit(event *i3.WindowEvent) error
	String() string
}

// Split based on width (W) and height (H)
type SmartSplitStrategy struct{}

func (strat *SmartSplitStrategy) SetUpNextSplit(event *i3.WindowEvent) error {
	command := "split h"
	if event.Container.Rect.Width < event.Container.Rect.Height {
		command = "split v"
	}
	_, err := i3.RunCommand(command)

	return err
}

func (strat *SmartSplitStrategy) String() string {
	return "smart"
}

func NewSmartSplitStrategy() *SmartSplitStrategy {
	return &SmartSplitStrategy{}
}

// Always split vertically
type VerticalSplitStrategy struct{}

func (strat *VerticalSplitStrategy) SetUpNextSplit(event *i3.WindowEvent) error {
	command := "split v"
	_, err := i3.RunCommand(command)

	return err
}

func (strat *VerticalSplitStrategy) String() string {
	return "vertical"
}

func NewVerticalSplitStrategy() *VerticalSplitStrategy {
	return &VerticalSplitStrategy{}
}

// Always split vertically
type HorizontalSplitStrategy struct{}

func (strat *HorizontalSplitStrategy) SetUpNextSplit(event *i3.WindowEvent) error {
	command := "split h"
	_, err := i3.RunCommand(command)

	return err
}

func (strat *HorizontalSplitStrategy) String() string {
	return "horizontal"
}

func NewHorizontalSplitStrategy() *HorizontalSplitStrategy {
	return &HorizontalSplitStrategy{}
}

// Split handling is off (same as default i3)
type OffSplitStrategy struct{}

func (strat *OffSplitStrategy) SetUpNextSplit(event *i3.WindowEvent) error {
	return nil
}

func (strat *OffSplitStrategy) String() string {
	return "off"
}

func NewOffSplitStrategy() *OffSplitStrategy {
	return &OffSplitStrategy{}
}
