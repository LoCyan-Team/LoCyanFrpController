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
        dir_name="LoCyanFrpController_${version}_${os}_${arch}"
        controller_path="./packages/LoCyanFrpController_${version}_${os}_${arch}"

        if [ "x${os}" = x"windows" ]; then
            if [ ! -f "./LoCyanFrpController_${os}_${arch}.exe" ]; then
                continue
            fi
            mkdir ${controller_path}
            mv ./LoCyanFrpController_${os}_${arch}.exe ${controller_path}/LoCyanFrpController.exe
        else
            if [ ! -f "./LoCyanFrpController_${os}_${arch}" ]; then
                continue
            fi
            mkdir ${controller_path}
            mv ./LoCyanFrpController_${os}_${arch} ${controller_path}/LoCyanFrpController
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