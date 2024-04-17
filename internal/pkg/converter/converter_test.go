package converter

import (
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestLogsJSONToXLS(t *testing.T) {
	testIn := `{"time":"2024-04-17T20:53:27.5487316+02:00","level":"INFO","msg":"setup storage"}
{"time":"2024-04-17T20:53:27.5495476+02:00","level":"INFO","msg":"setup services"}
{"time":"2024-04-17T20:53:27.5501848+02:00","level":"INFO","msg":"run server"}
{"time":"2024-04-17T20:54:57.3073816+02:00","level":"INFO","msg":"req","URL":"/student/grades","Method":"GET","Remote":"127.0.0.1:55419","body":""}
{"time":"2024-04-17T20:54:57.7757167+02:00","level":"INFO","msg":"req","URL":"/grades/by_month_and_user","Method":"POST","Remote":"127.0.0.1:55420","body":"{\"month\":9,\"course\":1}"}
{"time":"2024-04-17T21:08:11.4402727+02:00","level":"INFO","msg":"req","URL":"/grades/res_by_user","Method":"POST","Remote":"127.0.0.1:56205","body":"{\"course\":1}"}
{"time":"2024-04-17T21:10:14.4967917+02:00","level":"INFO","msg":"req","URL":"/student/grades","Method":"GET","Remote":"127.0.0.1:56205","body":""}
{"time":"2024-04-17T21:10:16.294703+02:00","level":"INFO","msg":"req","URL":"/grades/res_by_user","Method":"POST","Remote":"127.0.0.1:56296","body":"{\"course\":1}"}
{"time":"2024-04-17T21:10:18.6374548+02:00","level":"INFO","msg":"req","URL":"/grades/by_month_and_user","Method":"POST","Remote":"127.0.0.1:56298","body":"{\"month\":9,\"course\":1}"}
{"time":"2024-04-17T21:10:21.8032966+02:00","level":"INFO","msg":"req","URL":"/grades/res_by_user","Method":"POST","Remote":"127.0.0.1:56298","body":"{\"course\":1}"}
{"time":"2024-04-17T21:10:23.6615726+02:00","level":"INFO","msg":"req","URL":"/grades/by_month_and_user","Method":"POST","Remote":"127.0.0.1:56298","body":"{\"month\":12,\"course\":1}"}
{"time":"2024-04-17T21:10:25.8117941+02:00","level":"INFO","msg":"req","URL":"/grades/by_month_and_user","Method":"POST","Remote":"127.0.0.1:56298","body":"{\"month\":4,\"course\":1}"}
{"time":"2024-04-17T21:10:37.6479712+02:00","level":"INFO","msg":"req","URL":"/grades/res_by_user","Method":"POST","Remote":"127.0.0.1:56298","body":"{\"course\":1}"}
{"time":"2024-04-17T21:10:39.1102395+02:00","level":"INFO","msg":"req","URL":"/grades/by_month_and_user","Method":"POST","Remote":"127.0.0.1:56298","body":"{\"month\":9,\"course\":1}"}
{"time":"2024-04-17T22:13:33.3300123+02:00","level":"INFO","msg":"req","URL":"/grades/by_month_and_user","Method":"POST","Remote":"127.0.0.1:58616","body":"{\"month\":1,\"course\":5}"}
{"time":"2024-04-17T22:13:39.7233054+02:00","level":"INFO","msg":"req","URL":"/student/grades","Method":"GET","Remote":"127.0.0.1:58616","body":""}
{"time":"2024-04-17T22:13:39.777759+02:00","level":"INFO","msg":"req","URL":"/grades/by_month_and_user","Method":"POST","Remote":"127.0.0.1:58616","body":"{\"month\":1,\"course\":5}"}
{"time":"2024-04-17T22:13:41.8670161+02:00","level":"INFO","msg":"req","URL":"/grades/by_month_and_user","Method":"POST","Remote":"127.0.0.1:58616","body":"{\"month\":9,\"course\":1}"}
{"time":"2024-04-17T22:13:49.0764137+02:00","level":"INFO","msg":"req","URL":"/admin/users","Method":"GET","Remote":"127.0.0.1:58616","body":""}
{"time":"2024-04-17T22:13:49.1534653+02:00","level":"INFO","msg":"req","URL":"/users/get_all","Method":"POST","Remote":"127.0.0.1:58616","body":""}
{"time":"2024-04-17T22:13:50.5938942+02:00","level":"INFO","msg":"req","URL":"/admin/grades","Method":"GET","Remote":"127.0.0.1:58616","body":""}
{"time":"2024-04-17T22:13:50.8193495+02:00","level":"INFO","msg":"req","URL":"/grades/by_month_and_subject","Method":"POST","Remote":"127.0.0.1:58616","body":"{\"month\": 9, \"course\": 1, \"subject_id\": 15}"}
{"time":"2024-04-17T22:14:09.0816968+02:00","level":"INFO","msg":"req","URL":"/grades/by_month_and_subject","Method":"POST","Remote":"127.0.0.1:58616","body":"{\"month\": 3, \"course\": 5, \"subject_id\": 2}"}
{"time":"2024-04-17T22:14:09.1080889+02:00","level":"WARN","msg":"/httpServer/grades.go:99 > /subjects/subjects.go:35 > /mysql/subjects.go:30 > sql: no rows in result set"}`

	out, err := LogsJSONToXLS(testIn)
	require.NoError(t, err)

	fileData, err := io.ReadAll(out)
	require.NoError(t, err)

	err = os.WriteFile("test_logs.xlsx", fileData, 600)
	require.NoError(t, err)
}
