package slave

import (
  "fmt"
  "sort"
  "net"
  "os"
  "context"
  "encoding/json"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/containerd/containerd"
)

const (
  MASTER = "localhost:27001"
  importDir = "./slave-vol/"
  containerdSock = "/var/snap/microk8s/common/run/containerd.sock"
  containerdNameSpace = "cacis"
)

func Main() {
  slave()
}

func slave() {
  componentsList := recieveComponentsList()
  //fmt.Println(componentsList)
  // componentsList := map[string]string {
  //   "cni.img"             : "docker.io/calico/cni:v3.13.2",
  //   "pause.img"           : "k8s.gcr.io/pause:3.1",
  //   "kube-controllers.img": "docker.io/calico/kube-controllers:v3.13.2",
  //   "pod2daemon.img"      : "docker.io/calico/pod2daemon-flexvol:v3.13.2",
  //   "node.img"            : "docker.io/calico/node:v3.13.2",
  //   "coredns.img"         : "docker.io/coredns/coredns:1.8.0",
  //   "metrics-server.img"  : "k8s.gcr.io/metrics-server-arm64:v0.3.6",
  //   "dashboard.img"       : "docker.io/kubernetesui/dashboard:v2.0.0",
  // }
  sortedExportFileName := sortKeys(componentsList)
  recieveImg(sortedExportFileName)
  //importAllImg(componentsList)
}

func recieveComponentsList() map[string]string {
  fmt.Println("Debug: [start] RECIEVE COMPONENTS LIST")
  // Socket
  conn, err := net.Dial("tcp", MASTER)
  cacis.Error(err)
  defer conn.Close()

  // Request Image
  fmt.Println("Debug: Request Components List")
  cLayer := cacis.RequestComponentsList()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Println("Requested\n\n")


  // Recieve Packet
  fmt.Println("Debug: Recieve Packet")
  packet = make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(packet)
  cacis.Error(err)
  fmt.Printf("Debug: Read Packet HEADER. len: %d\n", packetLength)
  cLayer = cacis.Unmarshal(packet)
  //fmt.Println(cLayer)

  fmt.Println("Debug: Read Packet PAYLOAD")
  //fmt.Println(cLayer)
  cLayer.Payload = loadPayload(conn, cLayer.Length)

  fmt.Printf("\rCompleted  %d\n", len(cLayer.Payload))

  var tmp map[string]string
  err = json.Unmarshal(cLayer.Payload, &tmp)
  cacis.Error(err)

  fmt.Println("Debug: [end] RECIEVE COMPONENTS LIST")
  return tmp
}

func recieveImg(s []string) {
  fmt.Println("Debug: [start] RECIEVE COMPONENT IMAGES")
  // Socket
  conn, err := net.Dial("tcp", MASTER)
  cacis.Error(err)
  defer conn.Close()

  // Request Image
  fmt.Println("Debug: Request Components Image")
  cLayer := cacis.RequestImage()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Println("Requested\n\n")

  for _, fileName := range s {
    fmt.Printf("Debug: Recieve file '%s'\n", fileName)
    packet := make([]byte, cacis.CacisLayerSize)
    packetLength, err := conn.Read(packet)
    cacis.Error(err)
    fmt.Printf("Debug: Read Packet HEADER. len: %d\n", packetLength)
    cLayer = cacis.Unmarshal(packet)

    // Recieve Packet PAYLOAD
    fmt.Println("Debug: Read Packet PAYLOAD")
    cLayer.Payload = loadPayload(conn, cLayer.Length)

    fmt.Printf("\rCompleted  %d\n", len(cLayer.Payload))

    fmt.Printf("Debug: Write file '%s'\n\n", fileName)
    // File
    filePath := importDir + fileName
    //filePath := "./hoge1.txt"
    file , err := os.Create(filePath)
    cacis.Error(err)

    file.Write(cLayer.Payload)
  }
  fmt.Println("Debug: [end] RECIEVE COMPONENT IMAGES")
}

func importImg(imageName, filePath string) {
  fmt.Println("Importing " + imageName + " from " + filePath + "...")

  ctx := context.Background()
  client, err := containerd.New(containerdSock, containerd.WithDefaultNamespace(containerdNameSpace))
  defer client.Close()
  cacis.Error(err)

  f, err := os.Open(filePath)
  defer f.Close()
  cacis.Error(err)

  opts := []containerd.ImportOpt{
    containerd.WithIndexName(imageName),
    //containerd.WithAllPlatforms(true),
  }
  client.Import(ctx, f, opts...)
  cacis.Error(err)
  fmt.Println("Imported")
}

func importAllImg(m map[string]string) {
  fmt.Printf("Import %d images for Kubernetes Components", len(m))
  for importFile, imageRef := range m {
    fmt.Printf("%s : %s\n", importDir + importFile, imageRef)
    fmt.Println("start")
    importImg(imageRef, importDir + importFile)
    fmt.Println("end\n")
    }
}

func install_microk8s() {
}

func loadPayload(conn net.Conn, targetBytes uint64) []byte {
  packet := []byte{}
  recievedBytes := 0

  for len(packet) < int(targetBytes){
    buf := make([]byte, targetBytes - uint64(recievedBytes))
    packetLength, err := conn.Read(buf)
    cacis.Error(err)
    recievedBytes += packetLength
    packet = append(packet, buf[:packetLength]...)
    fmt.Printf("\rDebug: recieving...")
    fmt.Printf("\rCompleted  %d  of %d", len(packet), int(targetBytes))
  }
  return packet
}

func sortKeys(m map[string]string) []string {
  ///sort
  sorted := make([]string, len(m))
  index := 0
  for key := range m {
        sorted[index] = key
        index++
    }
    sort.Strings(sorted)
  /*
  for _, exportFile := range exportFileNameSort {
    fmt.Printf("%-20s : %s\n", exportFile, componentsList[exportFile])
  }
  */
  return sorted
}
