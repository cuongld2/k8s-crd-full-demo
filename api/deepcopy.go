package api

import "k8s.io/apimachinery/pkg/runtime"

func (in *Database) DeepCopyInto(out *Database) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = DatabaseSpec{
		DbName:      in.Spec.DbName,
		Description: in.Spec.Description,
		Total:       in.Spec.Total,
		Available:   in.Spec.Available,
		DbType:      in.Spec.DbType,
		Tags:        in.Spec.Tags,
	}
}

func (in *Database) DeepCopyObject() runtime.Object {
	out := Database{}
	in.DeepCopyInto(&out)

	return &out
}

func (in *DatabaseList) DeepCopyObject() runtime.Object {
	out := DatabaseList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]Database, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
