# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

version=0.0.1
echo "build version: $version"

# cross_compiles
make -f ./Makefile.cross-compiles

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin freebsd'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle riscv64'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        dir_name="controller_${version}_${os}_${arch}"
        controller_path="./packages/controller_${version}_${os}_${arch}"

        if [ "x${os}" = x"windows" ]; then
            if [ ! -f "./controller_${os}_${arch}.exe" ]; then
                continue
            fi
            mkdir ${controller_path}
            mv ./controller_${os}_${arch}.exe ${controller_path}/controller.exe
        else
            if [ ! -f "./controller_${os}_${arch}" ]; then
                continue
            fi
            mkdir ${controller_path}
            mv ./controller_${os}_${arch} ${controller_path}/controller
        fi  

        # packages
        cd ./packages
        if [ "x${os}" = x"windows" ]; then
            zip -rq ${dir_name}.zip ${dir_name}
        else
            tar -zcf ${dir_name}.tar.gz ${dir_name}
        fi  
        cd ..
        rm -rf ${controller_path}
    done
done

cd -