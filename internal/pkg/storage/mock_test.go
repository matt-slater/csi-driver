package storage_test

import (
	"csi-driver/internal/pkg/storage"
	"reflect"
	"testing"
)

func TestMockStorage_RemoveVolume(t *testing.T) {
	t.Parallel()

	type fields struct {
		ShouldErr bool
		Path      string
		Volumes   []string
	}

	type args struct {
		id string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "delete no err",
			fields: fields{},
			args: args{
				id: "test-id",
			},
		},
		{
			name: "delete err",
			fields: fields{
				ShouldErr: true,
			},
			args: args{
				id: "test-id",
			},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &storage.MockStorage{
				ShouldErr: testCase.fields.ShouldErr,
				Path:      testCase.fields.Path,
				Volumes:   testCase.fields.Volumes,
			}

			if err := mockStorage.RemoveVolume(testCase.args.id); (err != nil) != testCase.wantErr {
				t.Errorf("MockStorage.RemoveVolume() error = %v, wantErr %v", err, testCase.wantErr)
			}
		})
	}
}

func TestMockStorage_ListVolumes(t *testing.T) {
	t.Parallel()

	type fields struct {
		ShouldErr bool
		Path      string
		Volumes   []string
	}

	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name: "list volumes no err",
			fields: fields{
				Volumes: []string{"a", "b"},
			},
			want: []string{"a", "b"},
		},
		{
			name: "list volumes err",
			fields: fields{
				ShouldErr: true,
				Volumes:   []string{"a", "b"},
			},
			wantErr: true,
		},
	}
	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &storage.MockStorage{
				ShouldErr: testCase.fields.ShouldErr,
				Path:      testCase.fields.Path,
				Volumes:   testCase.fields.Volumes,
			}

			got, err := mockStorage.ListVolumes()
			if (err != nil) != testCase.wantErr {
				t.Errorf("MockStorage.ListVolumes() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("MockStorage.ListVolumes() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestMockStorage_PathForVolume(t *testing.T) {
	t.Parallel()

	type fields struct {
		ShouldErr bool
		Path      string
		Volumes   []string
	}

	type args struct {
		id string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "path for volume no err",
			fields: fields{
				Path: "/dev/test",
			},
			args: args{
				id: "id",
			},
			want: "/dev/test",
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &storage.MockStorage{
				ShouldErr: testCase.fields.ShouldErr,
				Path:      testCase.fields.Path,
				Volumes:   testCase.fields.Volumes,
			}

			if got := mockStorage.PathForVolume(testCase.args.id); got != testCase.want {
				t.Errorf("MockStorage.PathForVolume() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestMockStorage_WriteVolume(t *testing.T) {
	t.Parallel()

	type fields struct {
		ShouldErr bool
		Path      string
		Volumes   []string
	}

	type args struct {
		id   string
		vCtx map[string]string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "write volume no err",
			fields: fields{
				ShouldErr: false,
			},
			args:    args{},
			want:    true,
			wantErr: false,
		},
		{
			name: "write volume err",
			fields: fields{
				ShouldErr: true,
			},
			args:    args{},
			want:    false,
			wantErr: true,
		},
	}
	for _, testCase := range tests {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &storage.MockStorage{
				ShouldErr: testCase.fields.ShouldErr,
				Path:      testCase.fields.Path,
				Volumes:   testCase.fields.Volumes,
			}

			got, err := mockStorage.WriteVolume(testCase.args.id, testCase.args.vCtx)
			if (err != nil) != testCase.wantErr {
				t.Errorf("MockStorage.WriteVolume() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if got != testCase.want {
				t.Errorf("MockStorage.WriteVolume() = %v, want %v", got, testCase.want)
			}
		})
	}
}
