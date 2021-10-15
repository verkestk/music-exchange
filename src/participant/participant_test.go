package participant

import (
	"testing"
)

func Test_IsCompatible(t *testing.T) {
	woody := &Participant{Platforms: []string{}}
	buzz := &Participant{Platforms: []string{"Spotify"}}
	rex := &Participant{Platforms: []string{"Pandora"}}
	hamm := &Participant{Platforms: []string{"Pandora", "YouTube"}}
	jessie := &Participant{Platforms: []string{"Spotify", "AmazonMusic"}}

	if woody.IsCompatible(buzz) {
		t.Errorf("woody is compatible with buzz - expected incompatible")
	}

	if woody.IsCompatible(rex) {
		t.Errorf("woody is compatible with rex - expected incompatible")
	}

	if woody.IsCompatible(hamm) {
		t.Errorf("woody is compatible with hamm - expected incompatible")
	}

	if woody.IsCompatible(jessie) {
		t.Errorf("woody is compatible with jessie - expected incompatible")
	}

	if buzz.IsCompatible(woody) {
		t.Errorf("buzz is compatible with woody - expected incompatible")
	}

	if buzz.IsCompatible(rex) {
		t.Errorf("buzz is compatible with rex - expected incompatible")
	}

	if buzz.IsCompatible(hamm) {
		t.Errorf("buzz is compatible with hamm - expected incompatible")
	}

	if !buzz.IsCompatible(jessie) {
		t.Errorf("buzz is incompatible with jessie - expected compatible")
	}

	if rex.IsCompatible(woody) {
		t.Errorf("rex is compatible with woody - expected incompatible")
	}

	if rex.IsCompatible(buzz) {
		t.Errorf("rex is compatible with buzz - expected incompatible")
	}

	if !rex.IsCompatible(hamm) {
		t.Errorf("rex is incompatible with hamm - expected compatible")
	}

	if rex.IsCompatible(jessie) {
		t.Errorf("rex is compatible with jessie - expected incompatible")
	}

	if hamm.IsCompatible(woody) {
		t.Errorf("hamm is compatible with woody - expected incompatible")
	}

	if hamm.IsCompatible(buzz) {
		t.Errorf("hamm is compatible with buzz - expected incompatible")
	}

	if !hamm.IsCompatible(rex) {
		t.Errorf("hamm is incompatible with rex - expected compatible")
	}

	if hamm.IsCompatible(jessie) {
		t.Errorf("hamm is compatible with jessie - expected incompatible")
	}

	if jessie.IsCompatible(woody) {
		t.Errorf("jessie is compatible with woody - expected incompatible")
	}

	if !jessie.IsCompatible(buzz) {
		t.Errorf("jessie is incompatible with buzz - expected compatible")
	}

	if jessie.IsCompatible(rex) {
		t.Errorf("jessie is compatible with rex - expected incompatible")
	}

	if jessie.IsCompatible(hamm) {
		t.Errorf("jessie is compatible with hamm - expected incompatible")
	}
}
