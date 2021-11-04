package trace

//import (
//	"context"
//	"fmt"
//	"github.com/gmhafiz/go8/internal/domain/book"
//	"github.com/gmhafiz/go8/internal/models"
//	"github.com/opentracing/opentracing-go"
//	olog "github.com/opentracing/opentracing-go/log"
//	"time"
//)
//
//type Instrumented struct {
//	book.UseCase
//
//	Tracer opentracing.Tracer
//	Ctx    context.Context
//}
//
//var _ book.UseCase = (*Instrumented)(nil)
//
//type mT struct {
//	Tracer opentracing.Tracer
//	Ctx    context.Context
//}
//
//func NewMockTracer(ctx context.Context, tracer opentracing.Tracer) *Instrumented {
//
//	return &Instrumented{
//		Ctx:    ctx,
//		Tracer: tracer,
//	}
//}
//
//func (i *Instrumented) List(ctx context.Context, f *book.Filter) ([]*models.Book, error) {
//	span := opentracing.SpanFromContext(ctx)
//
//	span = i.Tracer.StartSpan("book_list")
//	span.SetTag(fmt.Sprintf("%s-called", "book_list"), time.Now())
//	span.LogFields(
//		olog.String("event", "book_list:get called"),
//	)
//	//defer span.Finish()
//
//	return i.UseCase.List(ctx, f)
//}
