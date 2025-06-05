package main

import (
	"reflect"
	"testing"

	"github.com/brvoisin/planetwars"
)

func Test_starterBot_FindPath(t *testing.T) {
	type args struct {
		pwMap  planetwars.Map
		source planetwars.PlanetID
		target planetwars.PlanetID
	}
	tests := []struct {
		name string
		args args
		want Path
	}{
		{
			name: "2 planets, direct path",
			args: args{
				pwMap: planetwars.Map{
					Planets: []planetwars.Planet{
						{ID: 0, Owner: planetwars.Myself},
						{ID: 1, Owner: planetwars.Opponent},
					},
				},
				source: 0,
				target: 1,
			},
			want: Path{0, 1},
		},
		{
			name: "Choose shortest path",
			args: args{
				pwMap: planetwars.Map{
					Planets: []planetwars.Planet{
						{ID: 0, Owner: planetwars.Myself, Position: planetwars.Point{X: 0, Y: 0}},
						{ID: 1, Owner: planetwars.Opponent, Position: planetwars.Point{X: 0, Y: 5}},
						{ID: 2, Owner: planetwars.Opponent, Position: planetwars.Point{X: 5, Y: 5}},
						{ID: 3, Owner: planetwars.Opponent, Position: planetwars.Point{X: 0, Y: 10}},
					},
				},
				source: 0,
				target: 3,
			},
			want: Path{0, 1, 3},
		},
		{
			name: "Choose shortest long path",
			args: args{
				pwMap: planetwars.Map{
					Planets: []planetwars.Planet{
						{ID: 0, Owner: planetwars.Myself, Position: planetwars.Point{X: 0, Y: 0}},
						{ID: 1, Owner: planetwars.Opponent, Position: planetwars.Point{X: 0, Y: 5}},
						{ID: 2, Owner: planetwars.Opponent, Position: planetwars.Point{X: 5, Y: 5}},
						{ID: 3, Owner: planetwars.Opponent, Position: planetwars.Point{X: 0, Y: 10}},
						{ID: 4, Owner: planetwars.Opponent, Position: planetwars.Point{X: 1, Y: 2}},
						{ID: 5, Owner: planetwars.Opponent, Position: planetwars.Point{X: 2, Y: 7}},
					},
				},
				source: 0,
				target: 3,
			},
			want: Path{0, 4, 1, 5, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &starterBot{}
			if got := b.FindPath(tt.args.pwMap, tt.args.source, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("starterBot.FindPath() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_starterBot_makeOrders(t *testing.T) {
	type fields struct {
		path Path
	}
	type args struct {
		pwMap planetwars.Map
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []planetwars.Order
	}{
		{
			name: "2 planets, one order with growth - 1",
			fields: fields{
				path: Path{0, 1},
			},
			args: args{
				pwMap: planetwars.Map{
					Planets: []planetwars.Planet{
						{ID: 0, Owner: planetwars.Myself, Ships: 10, Growth: 5},
						{ID: 1, Owner: planetwars.Opponent},
					},
				},
			},
			want: []planetwars.Order{
				{Source: 0, Dest: 1, Ships: 4},
			},
		},
		{
			name: "2 planets, one order with max ships - 1",
			fields: fields{
				path: Path{0, 1},
			},
			args: args{
				pwMap: planetwars.Map{
					Planets: []planetwars.Planet{
						{ID: 0, Owner: planetwars.Myself, Ships: 4, Growth: 5},
						{ID: 1, Owner: planetwars.Opponent},
					},
				},
			},
			want: []planetwars.Order{
				{Source: 0, Dest: 1, Ships: 3},
			},
		},
		{
			name: "3 planets, two orders with cumulative growth",
			fields: fields{
				path: Path{0, 1, 2},
			},
			args: args{
				pwMap: planetwars.Map{
					Planets: []planetwars.Planet{
						{ID: 0, Owner: planetwars.Myself, Ships: 10, Growth: 5},
						{ID: 1, Owner: planetwars.Myself, Ships: 10, Growth: 3},
						{ID: 2, Owner: planetwars.Opponent},
					},
				},
			},
			want: []planetwars.Order{
				{Source: 0, Dest: 1, Ships: 4},
				{Source: 1, Dest: 2, Ships: 6},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &starterBot{
				path: tt.fields.path,
			}
			if got := b.makeOrders(tt.args.pwMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("starterBot.makeOrders() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
